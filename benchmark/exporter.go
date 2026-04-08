// exporter.go exporta os resultados do benchmark para arquivos CSV.
//
// Dois arquivos são gerados:
//  - benchmark_results.csv  → uma linha por execução individual
//  - benchmark_stats.csv    → estatísticas agregadas por algoritmo × dataset
//
// O formato é compatível com Excel, Google Sheets, pandas e qualquer
// ferramenta de visualização que consuma CSV padrão.
package benchmark

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// ExportResults escreve uma linha por execução em path.
// Cabeçalho: Algorithm, DataType, Run, InputSize, DurationNs, DurationMs,
//            Comparisons, Swaps, MemAllocKB
func ExportResults(results []Result, path string) error {
	if err := os.MkdirAll(dirOf(path), 0755); err != nil {
		return fmt.Errorf("exporter: criar diretório: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("exporter: criar arquivo: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Cabeçalho
	header := []string{
		"Algorithm", "DataType", "Run", "InputSize",
		"DurationNs", "DurationMs",
		"Comparisons", "Swaps", "MemAllocKB",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	// Dados
	for _, r := range results {
		row := []string{
			r.Algorithm,
			r.DataType,
			strconv.Itoa(r.Run),
			strconv.Itoa(r.InputSize),
			strconv.FormatInt(r.Duration.Nanoseconds(), 10),
			formatFloat(r.DurationMs()),
			strconv.FormatInt(r.Comparisons, 10),
			strconv.FormatInt(r.Swaps, 10),
			formatFloat(r.MemAllocKB),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

// ExportStats escreve as estatísticas agregadas em path.
// Cabeçalho: Algorithm, DataType, Runs, MinMs, MaxMs, MeanMs, StdDevMs, MedianMs
func ExportStats(stats []Stats, path string) error {
	if err := os.MkdirAll(dirOf(path), 0755); err != nil {
		return fmt.Errorf("exporter: criar diretório: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("exporter: criar arquivo: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	header := []string{
		"Algorithm", "DataType", "Runs",
		"MinMs", "MaxMs", "MeanMs", "StdDevMs", "MedianMs",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	for _, s := range stats {
		row := []string{
			s.Algorithm,
			s.DataType,
			strconv.Itoa(s.Runs),
			formatFloat(s.MinMs),
			formatFloat(s.MaxMs),
			formatFloat(s.MeanMs),
			formatFloat(s.StdDevMs),
			formatFloat(s.MedianMs),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return w.Error()
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────

// formatFloat formata um float64 com 3 casas decimais, usando ponto como separador.
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 3, 64)
}

// dirOf retorna o diretório pai de um caminho de arquivo.
func dirOf(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}
