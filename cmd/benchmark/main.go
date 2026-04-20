// main.go — ponto de entrada do benchmark de algoritmos de ordenação.
//
// Uso:
//
//	go run cmd/benchmark/main.go
//
// Saída:
//
//	output/benchmark_results.csv  → execuções individuais
//	output/benchmark_stats.csv    → estatísticas agregadas
package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"pesquisa-ordenacao/benchmark"
	"pesquisa-ordenacao/data"
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════════════╗")
	fmt.Println("║     Benchmark de Algoritmos de Ordenação — Go        ║")
	fmt.Println("╚══════════════════════════════════════════════════════╝")
	fmt.Printf("  Seed fixa  : %d\n", data.Seed)
	fmt.Printf("  Tamanho    : %d elementos\n", data.Size)
	fmt.Printf("  Execuções  : %d por combinação\n\n", 3)

	// ── Registro dos algoritmos ────────────────────────────────────────────
	// Para adicionar um novo algoritmo: crie a struct em benchmark/algorithms.go
	// e adicione-a aqui. Nenhum outro arquivo precisa ser alterado.
	algorithms := []benchmark.SortAlgorithm{
		benchmark.BubbleSorter{},
		benchmark.SelectionSorter{},
		benchmark.InsertionSorter{},
		benchmark.ShellSorter{},
		benchmark.MergeSorter{},
		benchmark.QuickSorter{},
	}

	cfg := benchmark.DefaultConfig()

	// ── Aviso de tempo estimado ────────────────────────────────────────────
	// BubbleSort e InsertionSort são O(n²). Com 175.000 elementos, cada
	// execução no pior caso pode levar vários minutos.
	fmt.Println("⚠️  ATENÇÃO: algoritmos O(n²) com 175.000 elementos podem levar")
	fmt.Println("   vários minutos por execução. Total estimado: 20-60 minutos.")
	fmt.Println("   (ShellSort será muito mais rápido)")
	fmt.Println()

	wallStart := time.Now()

	// ── Execução do benchmark ──────────────────────────────────────────────
	results := benchmark.Run(algorithms, cfg)

	wallElapsed := time.Since(wallStart)

	// ── Estatísticas agregadas ─────────────────────────────────────────────
	stats := benchmark.ComputeStats(results)

	// ── Exportação para CSV ────────────────────────────────────────────────
	const (
		resultsPath = "output/benchmark_results.csv"
		statsPath   = "output/benchmark_stats.csv"
	)

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
