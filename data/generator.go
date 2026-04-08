// Package data fornece funções para geração de conjuntos de dados de teste
// para benchmarks de algoritmos de ordenação.
//
// Todos os datasets usam a mesma seed fixa (42) para garantir reprodutibilidade.
package data

import "math/rand"

// Constantes globais do benchmark.
const (
	// Seed fixa para reprodutibilidade absoluta entre execuções.
	Seed int64 = 42

	// Size é o tamanho padrão dos datasets (175.000 elementos).
	Size = 1000
)

// DataType representa o tipo de dataset gerado.
type DataType string

const (
	Sorted   DataType = "ordenado"
	Inverted DataType = "invertido"
	Random   DataType = "aleatorio"
)

// AllTypes retorna todos os tipos de dataset disponíveis, na ordem canônica.
func AllTypes() []DataType {
	return []DataType{Sorted, Inverted, Random}
}

// Generate cria e retorna um novo slice do tipo e tamanho especificados.
// Cada chamada retorna um slice independente; o caller é responsável por fazer
// cópias adicionais se necessário.
func Generate(size int, dt DataType) []int {
	switch dt {
	case Sorted:
		return generateSorted(size)
	case Inverted:
		return generateInverted(size)
	case Random:
		return generateRandom(size)
	default:
		panic("data.Generate: tipo desconhecido: " + string(dt))
	}
}

// generateSorted retorna [1, 2, 3, ..., size] — melhor caso para a maioria dos algoritmos.
func generateSorted(size int) []int {
	arr := make([]int, size)
	for i := range arr {
		arr[i] = i + 1
	}
	return arr
}

// generateInverted retorna [size, size-1, ..., 1] — pior caso para BubbleSort/InsertionSort.
func generateInverted(size int) []int {
	arr := make([]int, size)
	for i := range arr {
		arr[i] = size - i
	}
	return arr
}

// generateRandom retorna elementos aleatórios no intervalo [1, size*10]
// usando a seed fixa para reprodutibilidade.
func generateRandom(size int) []int {
	rng := rand.New(rand.NewSource(Seed))
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rng.Intn(size*10) + 1
	}
	return arr
}

// Clone retorna uma cópia profunda de um slice, segura para uso em benchmarks.
func Clone(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}
