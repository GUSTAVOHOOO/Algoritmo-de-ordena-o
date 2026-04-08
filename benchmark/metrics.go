// metrics.go define as estruturas de dados para armazenamento de resultados
// e as funções de cálculo de estatísticas agregadas (min, max, média, desvio padrão).
package benchmark

import (
	"math"
	"sort"
	"time"
)

// Result representa o resultado de uma única execução de benchmark.
type Result struct {
	Algorithm   string        // nome do algoritmo
	DataType    string        // tipo do dataset: ordenado | invertido | aleatorio
	Run         int           // índice da execução (1-based)
	Duration    time.Duration // tempo total de execução
	InputSize   int           // tamanho do slice de entrada
	Comparisons int64         // número de comparações elemento-a-elemento
	Swaps       int64         // número de trocas realizadas
	MemAllocKB  float64       // alocações de memória durante a execução (KB)
}

// DurationMs retorna a duração em milissegundos com precisão de 3 casas decimais.
func (r Result) DurationMs() float64 {
	return float64(r.Duration.Nanoseconds()) / 1e6
}

// ─────────────────────────────────────────────────────────────
// Estatísticas
// ─────────────────────────────────────────────────────────────

// Stats representa estatísticas agregadas para um par (algoritmo, tipo de dataset).
type Stats struct {
	Algorithm string
	DataType  string
	Runs      int
	MinMs     float64
	MaxMs     float64
	MeanMs    float64
	StdDevMs  float64
	MedianMs  float64
}

// ComputeStats calcula estatísticas por grupo (algoritmo × tipo de dataset).
func ComputeStats(results []Result) []Stats {
	type key struct{ algo, dt string }
	groups := make(map[key][]float64)

	for _, r := range results {
		k := key{r.Algorithm, r.DataType}
		groups[k] = append(groups[k], r.DurationMs())
	}

	var stats []Stats
	for k, durations := range groups {
		s := computeGroupStats(k.algo, k.dt, durations)
		stats = append(stats, s)
	}

	// Ordena por algoritmo e tipo para saída determinística
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Algorithm != stats[j].Algorithm {
			return stats[i].Algorithm < stats[j].Algorithm
		}
		return stats[i].DataType < stats[j].DataType
	})

	return stats
}

func computeGroupStats(algo, dt string, durations []float64) Stats {
	n := len(durations)
	if n == 0 {
		return Stats{Algorithm: algo, DataType: dt}
	}

	sorted := make([]float64, n)
	copy(sorted, durations)
	sort.Float64s(sorted)

	var total float64
	for _, d := range sorted {
		total += d
	}
	mean := total / float64(n)

	var variance float64
	for _, d := range sorted {
		diff := d - mean
		variance += diff * diff
	}
	variance /= float64(n)

	var median float64
	if n%2 == 0 {
		median = (sorted[n/2-1] + sorted[n/2]) / 2
	} else {
		median = sorted[n/2]
	}

	return Stats{
		Algorithm: algo,
		DataType:  dt,
		Runs:      n,
		MinMs:     sorted[0],
		MaxMs:     sorted[n-1],
		MeanMs:    mean,
		StdDevMs:  math.Sqrt(variance),
		MedianMs:  median,
	}
}
