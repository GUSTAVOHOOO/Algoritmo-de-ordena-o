// main.go — ponto de entrada do benchmark de algoritmos de ordenação.
//
// Uso:
//
//	go run cmd/benchmark/main.go
//
// Saída:
//
//	output/benchmark_<tipo>_<inputSize>_<iniciais>.csv
//	Ex: benchmark_results_175000_BISMRQ.csv
//	    benchmark_stats_175000_BISMRQ.csv
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"pesquisa-ordenacao/benchmark"
	"pesquisa-ordenacao/data"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║     Benchmark de Algoritmos de Ordenação — Go        ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")

	reader := bufio.NewReader(os.Stdin)

	// Captura do Tamanho do Dataset
	fmt.Printf("\n▶ Digite o número de elementos para o dataset (padrão: %d): ", data.Size)
	sizeInput, _ := reader.ReadString('\n')
	sizeInput = strings.TrimSpace(sizeInput)
	
	inputSize := data.Size
	if sizeInput != "" {
		if parsedSize, err := strconv.Atoi(sizeInput); err == nil && parsedSize > 0 {
			inputSize = parsedSize
		} else {
			fmt.Printf("  [!] Valor inválido. Usando padrão %d.\n", inputSize)
		}
	}

	// ── Registro de todos os algoritmos ────────────────────────────────────
	allAlgorithms := []benchmark.SortAlgorithm{
		benchmark.BubbleSorter{},
		benchmark.SelectionSorter{},
		benchmark.InsertionSorter{},
		benchmark.ShellSorter{},
		benchmark.MergeSorter{},
		benchmark.QuickSorter{},
		benchmark.HeapSorter{},
		benchmark.RadixSorter{},
	}

	// Captura das opções de Algoritmo
	fmt.Println("\n▶ Selecione os algoritmos para rodar:")
	fmt.Println("   0 - TODOS")
	for i, algo := range allAlgorithms {
		fmt.Printf("   %d - %s\n", i+1, algo.Name())
	}
	fmt.Print("  Escolha (ex: 0, ou '1,4,5') [padrão: 0]: ")
	algoInput, _ := reader.ReadString('\n')
	algoInput = strings.TrimSpace(algoInput)

	var selectedAlgorithms []benchmark.SortAlgorithm
	if algoInput == "" || algoInput == "0" {
		selectedAlgorithms = allAlgorithms
	} else {
		choices := strings.Split(algoInput, ",")
		for _, choice := range choices {
			idx, err := strconv.Atoi(strings.TrimSpace(choice))
			if err == nil && idx >= 1 && idx <= len(allAlgorithms) {
				selectedAlgorithms = append(selectedAlgorithms, allAlgorithms[idx-1])
			}
		}
		if len(selectedAlgorithms) == 0 {
			fmt.Println("  [!] Escolhas inválidas. Adotando TODOS como padrão.")
			selectedAlgorithms = allAlgorithms
		}
	}

	fmt.Println("\n========================================================")
	fmt.Printf("  Seed fixa  : %d\n", data.Seed)
	fmt.Printf("  Tamanho    : %d elementos\n", inputSize)
	fmt.Printf("  Algoritmos : %d selecionado(s)\n", len(selectedAlgorithms))
	fmt.Printf("  Execuções  : %d por combinação\n", 3)
	fmt.Println("========================================================")
	fmt.Println()

	cfg := benchmark.DefaultConfig()
	cfg.InputSize = inputSize

	// ── Aviso de tempo estimado ────────────────────────────────────────────
	// BubbleSort e InsertionSort são O(n²). Com entradas gigantes, pode
	// levar minutos.
	if inputSize >= 150000 {
		fmt.Printf("⚠️  ATENÇÃO: algoritmos O(n²) com %d elementos podem levar\n", inputSize)
		fmt.Println("   vários minutos por execução.")
		fmt.Println()
	}

	wallStart := time.Now()

	// ── Execução do benchmark ──────────────────────────────────────────────
	results := benchmark.Run(selectedAlgorithms, cfg)

	wallElapsed := time.Since(wallStart)

	// ── Estatísticas agregadas ─────────────────────────────────────────────
	stats := benchmark.ComputeStats(results)

	// ── Exportação para CSV ────────────────────────────────────────────────
	resultsPath := "output/" + benchmark.GenerateFilename("results", inputSize, selectedAlgorithms)
	statsPath := "output/" + benchmark.GenerateFilename("stats", inputSize, selectedAlgorithms)

	fmt.Println("💾 Exportando resultados...")

	if err := benchmark.ExportResults(results, resultsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao exportar resultados: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ %s\n", resultsPath)

	if err := benchmark.ExportStats(stats, statsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao exportar estatísticas: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  ✓ %s\n\n", statsPath)

	// ── Tabela de resumo no terminal ───────────────────────────────────────
	fmt.Println("📈 RESUMO (tempo médio em ms):")
	fmt.Println()

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "Algoritmo\tDataset\tMin (ms)\tMédia (ms)\tMax (ms)\tDesvPad (ms)\t")
	fmt.Fprintln(tw, "---------\t-------\t--------\t----------\t---------\t------------\t")

	for _, s := range stats {
		fmt.Fprintf(tw, "%s\t%s\t%.3f\t%.3f\t%.3f\t%.3f\t\n",
			s.Algorithm, s.DataType,
			s.MinMs, s.MeanMs, s.MaxMs, s.StdDevMs,
		)
	}
	tw.Flush()

	fmt.Printf("\n⏱  Tempo total de execução: %v\n", wallElapsed.Round(time.Second))
	fmt.Println("\n✅ Benchmark concluído.")
}
