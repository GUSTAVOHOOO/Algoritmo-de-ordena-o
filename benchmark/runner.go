// runner.go implementa o motor de execução do benchmark.
//
// Responsabilidades:
//  - Pré-gerar todos os datasets antes das medições
//  - Garantir cópia limpa dos dados por execução (sem efeito colateral)
//  - Forçar GC antes de cada medição para minimizar ruído
//  - Medir tempo e memória de forma isolada
//  - Emitir progresso no stderr durante a execução
package benchmark

import (
	"fmt"
	"runtime"
	"time"

	"pesquisa-ordenacao/data"
)

// Config controla o comportamento do runner.
type Config struct {
	Runs      int  // número de execuções por combinação algoritmo × dataset
	InputSize int  // tamanho do slice de entrada
	Verbose   bool // emite progresso durante a execução
}

// DefaultConfig retorna a configuração padrão conforme especificação.
func DefaultConfig() Config {
	return Config{
		Runs:      3,
		InputSize: data.Size,
		Verbose:   true,
	}
}

// Run executa o benchmark completo e retorna todos os resultados individuais.
//
// Fluxo:
//  1. Gera os 3 datasets uma única vez e mantém-os em memória como referência.
//  2. Para cada algoritmo × dataset × run:
//     a. Clona o dataset original (cópia limpa, sem modificar a referência).
//     b. Força o GC para eliminar ruído de alocações anteriores.
//     c. Captura MemStats antes e depois para medir alocações da execução.
//     d. Mede o tempo com time.Now() / time.Since().
//     e. Armazena o Result.
func Run(algorithms []SortAlgorithm, cfg Config) []Result {
	dataTypes := data.AllTypes()

	// ── 1. Pré-geração dos datasets de referência ──────────────────────────
	// Gerados uma única vez para garantir que todos os algoritmos recebam
	// exatamente os mesmos dados de entrada.
	if cfg.Verbose {
		fmt.Printf("\n📊 Gerando datasets (%d elementos cada)...\n", cfg.InputSize)
	}
	datasets := make(map[data.DataType][]int, len(dataTypes))
	for _, dt := range dataTypes {
		datasets[dt] = data.Generate(cfg.InputSize, dt)
		if cfg.Verbose {
			fmt.Printf("  ✓ %-12s gerado\n", string(dt))
		}
	}

	totalRuns := len(algorithms) * len(dataTypes) * cfg.Runs
	current := 0

	if cfg.Verbose {
		fmt.Printf("\n🚀 Iniciando benchmark: %d algoritmos × %d datasets × %d execuções = %d execuções totais\n\n",
			len(algorithms), len(dataTypes), cfg.Runs, totalRuns)
	}

	var results []Result

	// ── 2. Loop principal ──────────────────────────────────────────────────
	for _, algo := range algorithms {
		if cfg.Verbose {
			fmt.Printf("► %s\n", algo.Name())
		}

		for _, dt := range dataTypes {
			original := datasets[dt]

			for run := 1; run <= cfg.Runs; run++ {
				current++

				// a. Cópia limpa — evita que uma execução afete a próxima
				arr := data.Clone(original)

				// b. Força GC para limpar alocações anteriores
				runtime.GC()
				runtime.GC() // duas passagens para maior garantia

				// c. Memória antes
				var memBefore, memAfter runtime.MemStats
				runtime.ReadMemStats(&memBefore)

				// d. Medição de tempo — nenhum print ou I/O entre start e since
				start := time.Now()
				sr := algo.Sort(arr)
				elapsed := time.Since(start)

				// e. Memória depois
				runtime.ReadMemStats(&memAfter)

				// Alocações durante a execução em KB
				allocBytes := float64(memAfter.TotalAlloc-memBefore.TotalAlloc) / 1024.0

				// Usa o resultado para evitar que o compilador otimize a chamada
				_ = sr.Sorted

				results = append(results, Result{
					Algorithm:   algo.Name(),
					DataType:    string(dt),
					Run:         run,
					Duration:    elapsed,
					InputSize:   cfg.InputSize,
					Comparisons: sr.Comparisons,
					Swaps:       sr.Swaps,
					MemAllocKB:  allocBytes,
				})

				if cfg.Verbose {
					fmt.Printf("  [%2d/%2d] %-14s | %-10s | run %d → %10.3f ms | cmps: %15d | swaps: %13d\n",
						current, totalRuns,
						algo.Name(), string(dt), run,
						float64(elapsed.Nanoseconds())/1e6,
						sr.Comparisons, sr.Swaps,
					)
				}
			}
		}

		if cfg.Verbose {
			fmt.Println()
		}
	}

	return results
}
