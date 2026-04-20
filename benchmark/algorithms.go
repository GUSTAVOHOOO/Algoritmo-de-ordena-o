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

// ─────────────────────────────────────────────────────────────
// Merge Sort
// Complexidade: O(n log n) em todos os casos; Memória O(n)
// ─────────────────────────────────────────────────────────────

// MergeSorter implementa SortAlgorithm com contagem de operações.
type MergeSorter struct{}

func (MergeSorter) Name() string { return "MergeSort" }

func (MergeSorter) Sort(arr []int) SortResult {
	var cmp, swaps int64
	res := mergeSortInstrumented(arr, &cmp, &swaps)
	return SortResult{res, cmp, swaps}
}

func mergeSortInstrumented(arr []int, cmp, swaps *int64) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := mergeSortInstrumented(arr[:mid], cmp, swaps)
	right := mergeSortInstrumented(arr[mid:], cmp, swaps)
	return mergeInstrumented(left, right, cmp, swaps)
}

func mergeInstrumented(left, right []int, cmp, swaps *int64) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		*cmp++
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
		*swaps++
	}
	for i < len(left) {
		result = append(result, left[i])
		i++
		*swaps++
	}
	for j < len(right) {
		result = append(result, right[j])
		j++
		*swaps++
	}
	return result
}

// ─────────────────────────────────────────────────────────────
// Quick Sort
// Complexidade: O(n log n) médio; O(n²) pior caso
// ─────────────────────────────────────────────────────────────

// QuickSorter implementa SortAlgorithm com contagem de operações.
type QuickSorter struct{}

func (QuickSorter) Name() string { return "QuickSort" }

func (QuickSorter) Sort(arr []int) SortResult {
	var cmp, swaps int64
	quickSortInstrumented(arr, 0, len(arr)-1, &cmp, &swaps)
	return SortResult{arr, cmp, swaps}
}

func quickSortInstrumented(arr []int, low, high int, cmp, swaps *int64) {
	if low < high {
		p := partitionInstrumented(arr, low, high, cmp, swaps)
		quickSortInstrumented(arr, low, p, cmp, swaps)
		quickSortInstrumented(arr, p+1, high, cmp, swaps)
	}
}

func partitionInstrumented(arr []int, low, high int, cmp, swaps *int64) int {
	// Pivô mediano para mitigar O(n²) em entradas já ordenadas/invertidas
	pivot := arr[(low+high)/2]
	i := low - 1
	j := high + 1
	for {
		for {
			i++
			*cmp++
			if arr[i] >= pivot {
				break
			}
		}
		for {
			j--
			*cmp++
			if arr[j] <= pivot {
				break
			}
		}
		if i >= j {
			return j
		}
		arr[i], arr[j] = arr[j], arr[i]
		*swaps++
	}
}

// ─────────────────────────────────────────────────────────────
// Heap Sort
// Complexidade: O(n log n) — todos os casos; Memória O(1)
// ─────────────────────────────────────────────────────────────

// HeapSorter implementa SortAlgorithm com contagem de operações.
type HeapSorter struct{}

func (HeapSorter) Name() string { return "HeapSort" }

func (HeapSorter) Sort(arr []int) SortResult {
	var cmp, swaps int64
	n := len(arr)

	heapDown := func(i, size int) {
		for {
			largest := i
			l, r := 2*i+1, 2*i+2
			if l < size {
				cmp++
				if arr[l] > arr[largest] {
					largest = l
				}
			}
			if r < size {
				cmp++
				if arr[r] > arr[largest] {
					largest = r
				}
			}
			if largest == i {
				break
			}
			arr[i], arr[largest] = arr[largest], arr[i]
			swaps++
			i = largest
		}
	}

	for i := n/2 - 1; i >= 0; i-- {
		heapDown(i, n)
	}
	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		swaps++
		heapDown(0, i)
	}
	return SortResult{arr, cmp, swaps}
}

// ─────────────────────────────────────────────────────────────
// Radix Sort
// Complexidade: O(n·k) — todos os casos; k = dígitos do maior valor
// ─────────────────────────────────────────────────────────────

// RadixSorter implementa SortAlgorithm com contagem de operações.
// Comparações são zero (radix não compara elementos diretamente);
// Swaps contam cada escrita no buffer de saída.
type RadixSorter struct{}

func (RadixSorter) Name() string { return "RadixSort" }

func (RadixSorter) Sort(arr []int) SortResult {
	var swaps int64

	countSort := func(a []int, exp int) []int {
		var digits [10]int
		out := make([]int, len(a))
		for _, v := range a {
			digits[(v/exp)%10]++
		}
		for i := 1; i < 10; i++ {
			digits[i] += digits[i-1]
		}
		for i := len(a) - 1; i >= 0; i-- {
			idx := (a[i]/exp)%10
			digits[idx]--
			out[digits[idx]] = a[i]
			swaps++
		}
		return out
	}

	unsignedRadix := func(a []int) []int {
		if len(a) == 0 {
			return a
		}
		mx := a[0]
		for _, v := range a[1:] {
			if v > mx {
				mx = v
			}
		}
		for exp := 1; mx/exp > 0; exp *= 10 {
			a = countSort(a, exp)
		}
		return a
	}

	var neg, pos []int
	for _, v := range arr {
		if v < 0 {
			neg = append(neg, -v)
		} else {
			pos = append(pos, v)
		}
	}

	neg = unsignedRadix(neg)
	for i, j := 0, len(neg)-1; i <= j; i, j = i+1, j-1 {
		neg[i], neg[j] = -neg[j], -neg[i]
	}
	pos = unsignedRadix(pos)

	copy(arr, append(neg, pos...))
	return SortResult{arr, 0, swaps}
}
