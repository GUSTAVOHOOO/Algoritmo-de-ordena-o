// Package benchmark implementa algoritmos de ordenação instrumentados:
// cada implementação conta comparações e trocas além de ordenar o slice.
//
// Design: a interface SortAlgorithm desacopla o runner dos algoritmos concretos,
// permitindo adicionar novos algoritmos sem alterar nenhum outro arquivo.
package benchmark

// SortResult encapsula o resultado de uma execução de ordenação instrumentada.
type SortResult struct {
	Sorted      []int // slice ordenado (mesmo que o de entrada, modificado in-place)
	Comparisons int64 // número de comparações elemento-a-elemento
	Swaps       int64 // número de trocas (movimentações de elementos)
}

// SortAlgorithm é a interface comum que todo algoritmo deve implementar.
// Para adicionar um novo algoritmo:
//  1. Crie uma struct que implemente essa interface.
//  2. Adicione-a à lista em cmd/benchmark/main.go.
type SortAlgorithm interface {
	Name() string
	Sort(arr []int) SortResult
}

// ─────────────────────────────────────────────────────────────
// Bubble Sort
// Complexidade: O(n²) — pior e médio; O(n) — melhor (lista já ordenada)
// ─────────────────────────────────────────────────────────────

// BubbleSorter implementa SortAlgorithm com contagem de operações.
type BubbleSorter struct{}

func (BubbleSorter) Name() string { return "BubbleSort" }

func (BubbleSorter) Sort(arr []int) SortResult {
	var cmp, swaps int64
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(arr)-1; i++ {
			cmp++
			if arr[i+1] < arr[i] {
				arr[i+1], arr[i] = arr[i], arr[i+1]
				swapped = true
				swaps++
			}
		}
	}
	return SortResult{arr, cmp, swaps}
}

// ─────────────────────────────────────────────────────────────
// Selection Sort
// Complexidade: O(n²) — todos os casos; O(n) trocas no máximo
// ─────────────────────────────────────────────────────────────

// SelectionSorter implementa SortAlgorithm com contagem de operações.
type SelectionSorter struct{}

func (SelectionSorter) Name() string { return "SelectionSort" }

func (SelectionSorter) Sort(arr []int) SortResult {
	var cmp, swaps int64
	n := len(arr)
	for i := 0; i < n; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			cmp++
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		if minIdx != i {
			arr[i], arr[minIdx] = arr[minIdx], arr[i]
			swaps++
		}
	}
	return SortResult{arr, cmp, swaps}
}

// ─────────────────────────────────────────────────────────────
// Insertion Sort
// Complexidade: O(n²) — pior e médio; O(n) — melhor (lista já ordenada)
// ─────────────────────────────────────────────────────────────

// InsertionSorter implementa SortAlgorithm com contagem de operações.
type InsertionSorter struct{}

func (InsertionSorter) Name() string { return "InsertionSort" }

func (InsertionSorter) Sort(arr []int) SortResult {
	var cmp, swaps int64
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 {
			cmp++
			if arr[j] <= key {
				break
			}
			arr[j+1] = arr[j]
			swaps++
			j--
		}
		arr[j+1] = key
	}
	return SortResult{arr, cmp, swaps}
}

// ─────────────────────────────────────────────────────────────
// Shell Sort
// Complexidade: depende da sequência de gaps — geralmente O(n log²n)
// ─────────────────────────────────────────────────────────────

// ShellSorter implementa SortAlgorithm com contagem de operações.
type ShellSorter struct{}

func (ShellSorter) Name() string { return "ShellSort" }

func (ShellSorter) Sort(arr []int) SortResult {
	var cmp, swaps int64
	n := len(arr)
	for gap := n / 2; gap > 0; gap /= 2 {
		for i := gap; i < n; i++ {
			key := arr[i]
			j := i
			for j >= gap {
				cmp++
				if arr[j-gap] <= key {
					break
				}
				arr[j] = arr[j-gap]
				swaps++
				j -= gap
			}
			arr[j] = key
		}
	}
	return SortResult{arr, cmp, swaps}
}
