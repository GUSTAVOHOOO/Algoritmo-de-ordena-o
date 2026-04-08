# Benchmark de Algoritmos de Ordenação

> **Trabalho acadêmico** desenvolvido no âmbito da disciplina de **Pesquisa e Ordenação**  
> [Universidade Tecnológica Federal do Paraná — UTFPR](https://www.utfpr.edu.br/)

Este repositório contém uma infraestrutura de benchmarking científico para análise experimental do desempenho de algoritmos de ordenação clássicos. O projeto está em desenvolvimento contínuo: novos algoritmos serão incorporados ao longo da disciplina, e os resultados experimentais completos subsidiarão a elaboração de um **artigo científico** ao final do semestre.

---

## Sumário

- [Contexto e Motivação](#contexto-e-motivação)
- [Por que Go?](#por-que-go)
- [Algoritmos Implementados](#algoritmos-implementados)
- [Metodologia Experimental](#metodologia-experimental)
- [Estrutura do Repositório](#estrutura-do-repositório)
- [Tech Stack](#tech-stack)
- [Pré-requisitos](#pré-requisitos)
- [Como Executar](#como-executar)
- [Geração de Gráficos](#geração-de-gráficos)
- [Métricas Coletadas](#métricas-coletadas)
- [Resultados Parciais](#resultados-parciais)
- [Adicionando Novos Algoritmos](#adicionando-novos-algoritmos)
- [Próximas Etapas](#próximas-etapas)

---

## Contexto e Motivação

A análise de complexidade assintótica — expressa em notação **Big-O** — constitui um dos pilares teóricos da Ciência da Computação. Contudo, a complexidade assintótica descreve o comportamento *no limite*, abstraindo constantes, coeficientes e efeitos de hardware como localidade de cache, pipeline do processador e custo de alocação de memória.

Este projeto adota uma abordagem **empírica e reprodutível**:

1. Implementar os algoritmos com **instrumentação in-code** (contagem exata de comparações e trocas).
2. Executar cada algoritmo sobre três tipos de dataset — **ordenado**, **invertido** e **aleatório** — em múltiplos tamanhos de entrada.
3. Registrar métricas de tempo de execução, operações elementares e alocações de memória.
4. Visualizar e comparar os resultados para validar (ou contestar) as previsões teóricas em condições reais de hardware.

O objetivo final é produzir evidências experimentais que complementem a análise teórica e que possam ser reportadas em formato científico.

---

## Por que Go?

A escolha da linguagem **Go (Golang)** para este benchmark não é arbitrária. Ela responde a requisitos científicos precisos:

### 1. Compilação para código nativo (assembly)

Go é uma linguagem **compilada estaticamente**: o compilador `gc` gera código de máquina nativo (x86-64, ARM64, etc.) diretamente, sem intermediários como bytecode ou interpretação em tempo de execução. Isso significa que as operações de comparação e troca executam em instruções de assembly, sem overhead de JIT warm-up (como em Java/JVM) ou de interpretação (como em Python).

Para benchmarks de algoritmos, isso é fundamental: o tempo medido reflete o custo computacional real das operações, não artefatos do ambiente de execução.

### 2. Controle explícito do Garbage Collector

Go permite invocar o GC manualmente via `runtime.GC()`. O runner do benchmark faz **duas chamadas explícitas ao GC** antes de cada medição:

```go
runtime.GC()
runtime.GC() // duas passagens para maior garantia
```

Isso elimina o ruído causado por coletas de lixo de ciclos anteriores contaminando a medição atual — um problema crítico em linguagens com GC automático como Java ou C#.

### 3. Medição de tempo de alta resolução

`time.Since(start)` em Go usa o relógio monotônico do sistema operacional com resolução de **nanossegundos**, sem conversões de tipo ou arredondamentos intermediários. A medição é feita com o mínimo de código entre `start` e `elapsed`:

```go
start := time.Now()
sr := algo.Sort(arr)      // única instrução entre as marcas de tempo
elapsed := time.Since(start)
```

### 4. Ausência de custo oculto por abstrações

A interface `SortAlgorithm` em Go usa **dispatch estático via tabela de métodos** (similar ao vtable do C++), sem boxing, reflexão em tempo de execução ou overhead de closures. Os algoritmos operam diretamente sobre slices (`[]int`), que são apenas um ponteiro + comprimento + capacidade — sem wrapper objects.

### 5. Reprodutibilidade via seed fixa

Go permite fixar a seed do gerador pseudoaleatório:

```go
rng := rand.New(rand.NewSource(42))
```

Todos os runs produzem exatamente os mesmos datasets aleatórios, garantindo **reprodutibilidade absoluta** dos experimentos em qualquer máquina.

### Comparação com alternativas

| Critério | Go | Python | Java | C/C++ |
|---|:---:|:---:|:---:|:---:|
| Compilado para nativo | ✅ | ❌ | ❌ (JVM) | ✅ |
| Controle do GC | ✅ | ❌ | Parcial | N/A |
| Tempo de desenvolvimento | ✅ | ✅ | ❌ | ❌ |
| Sem overhead de JIT warm-up | ✅ | ✅ | ❌ | ✅ |
| Stdlib robusta para benchmarks | ✅ | Parcial | ✅ | ❌ |

Go oferece o equilíbrio ideal entre **precisão de medição** (próxima de C) e **produtividade de desenvolvimento** (próxima de Python).

---

## Algoritmos Implementados

### Fase atual

| Algoritmo | Complexidade (pior) | Complexidade (médio) | Complexidade (melhor) | Complexidade de espaço |
|---|:---:|:---:|:---:|:---:|
| **Bubble Sort** | O(n²) | O(n²) | **O(n)** | O(1) |
| **Selection Sort** | O(n²) | O(n²) | O(n²) | O(1) |
| **Insertion Sort** | O(n²) | O(n²) | **O(n)** | O(1) |
| **Shell Sort** | O(n log²n)* | O(n log²n)* | O(n log n) | O(1) |

> *Shell Sort com sequência de Knuth (gap = n/2, n/4, ..., 1). A complexidade exata depende da sequência de gaps utilizada.

### Algoritmos previstos para as próximas fases

- [ ] Merge Sort — O(n log n) estável, com custo de memória O(n)
- [ ] Quick Sort — O(n log n) médio, O(n²) pior caso
- [ ] Heap Sort — O(n log n) garantido, in-place
- [ ] Counting Sort — O(n + k), não comparativo
- [ ] Radix Sort — O(nk), não comparativo
- [ ] Tim Sort — híbrido Merge + Insertion, usado em Python e Java

---

## Metodologia Experimental

### Design do experimento

Cada combinação `(algoritmo, tipo de dataset, tamanho de entrada)` é executada **3 vezes** de forma independente, com as seguintes garantias:

1. **Isolamento de dados**: cada execução recebe uma cópia limpa (`data.Clone`) do dataset original, evitando que uma execução afete as seguintes (os algoritmos modificam o slice in-place).
2. **Limpeza do GC**: duas chamadas a `runtime.GC()` antes de cada medição.
3. **Medição mínima**: nenhuma instrução de I/O entre `time.Now()` e `time.Since()`.
4. **Coleta de memória**: `runtime.MemStats` capturado antes e depois de cada sort.

### Tipos de dataset

| Tipo | Descrição | Caso teórico |
|---|---|---|
| `ordenado` | `[1, 2, 3, ..., n]` | Melhor caso para Bubble e Insertion Sort |
| `invertido` | `[n, n-1, ..., 1]` | Pior caso para Bubble e Insertion Sort |
| `aleatorio` | Permutação aleatória com seed fixa 42 | Caso médio |

### Tamanhos de entrada avaliados

| Tamanho | Execuções totais por algoritmo |
|---|---|
| 175.000 elementos | 3 runs × 3 datasets = 9 |
| 250.000 elementos | 9 |
| 500.000 elementos | 9 |
| 750.000 elementos | 9 |
| 1.500.000 elementos | 9 |

### Estatísticas computadas

Para cada grupo `(algoritmo, dataset)`, são calculados:

- **Mínimo** — melhor tempo observado nas 3 execuções
- **Máximo** — pior tempo observado
- **Média aritmética** — estimador central
- **Desvio padrão** — variabilidade / ruído da medição
- **Mediana** — mais robusta a outliers que a média

---

## Estrutura do Repositório

```
.
├── benchmark/
│   ├── algorithms.go   # Implementações instrumentadas dos algoritmos
│   │                   # (interface SortAlgorithm + structs concretas)
│   ├── exporter.go     # Exportação de resultados e estatísticas para CSV
│   ├── metrics.go      # Structs Result e Stats; cálculo de estatísticas
│   └── runner.go       # Motor de benchmarking (geração, isolamento, medição)
│
├── cmd/
│   └── benchmark/
│       └── main.go     # Ponto de entrada: registra algoritmos e executa o benchmark
│
├── data/
│   └── generator.go    # Geração de datasets (ordenado, invertido, aleatório)
│                       # Seed fixa 42 para reprodutibilidade
│
├── scripts/
│   └── plot_results.py # Geração de gráficos e tabelas com matplotlib/seaborn
│
├── output/
│   ├── benchmark_results.csv        # Execuções individuais (linha por run)
│   ├── benchmark_results_750k.csv   # Idem para 750k elementos
│   ├── benchmark_results_1500k.csv  # Idem para 1.5M elementos
│   ├── benchmark_stats.csv          # Estatísticas agregadas
│   ├── benchmark_stats_750k.csv
│   ├── benchmark_stats_1500k.csv
│   └── plots_python/
│       ├── 175k/   → py_tempo_medio-175.png | py_tabela_resultados175.png | py_tabela_completa175.png
│       ├── 250k/   → py_tempo_medio250.png  | ...
│       ├── 500k/   → py_tempo_medio500.png  | ...
│       ├── 750k/   → py_tempo_medio750.png  | ...
│       └── 1500k/  → py_tempo_medio1500.png | ...
│
├── go.mod
├── go.sum
├── .gitignore
└── README.md
```

---

## Tech Stack

| Camada | Tecnologia | Versão | Papel |
|---|---|---|---|
| Benchmark | Go | 1.23+ | Implementação, instrumentação e medição |
| Visualização | Python | 3.10+ | Geração de gráficos e tabelas |
| Plots | matplotlib + seaborn | latest | Gráficos de barras com escala logarítmica |
| Tabelas | dataframe-image | latest | Exportação de DataFrames pandas para PNG |
| Dados | pandas | latest | Leitura e agregação dos CSVs |

---

## Pré-requisitos

### Go

- [Go 1.21+](https://go.dev/dl/) instalado e no `PATH`
- Verificar: `go version`

### Python

- Python 3.10+
- Biblioteca `pip` disponível
- Dependências:

```bash
pip install pandas matplotlib seaborn dataframe-image
```

---

## Como Executar

### 1. Clonar o repositório

```bash
git clone https://github.com/GUSTAVOHOOO/Algoritmo-de-ordena-o.git
cd "Algoritmo-de-ordena-o"
```

### 2. Configurar o tamanho do dataset

Edite a constante `Size` em `data/generator.go`:

```go
const (
    Seed int64 = 42
    Size       = 175_000  // altere aqui: 250_000, 500_000, etc.
)
```

### 3. Executar o benchmark

```bash
go run cmd/benchmark/main.go
```

**Saída esperada no terminal:**

```
╔══════════════════════════════════════════════════════╗
║     Benchmark de Algoritmos de Ordenação — Go        ║
╚══════════════════════════════════════════════════════╝
  Seed fixa  : 42
  Tamanho    : 175000 elementos
  Execuções  : 3 por combinação

📊 Gerando datasets (175000 elementos cada)...
  ✓ ordenado     gerado
  ✓ invertido    gerado
  ✓ aleatorio    gerado

🚀 Iniciando benchmark: 4 algoritmos × 3 datasets × 3 execuções = 36 execuções totais

► BubbleSort
  [ 1/36] BubbleSort     | ordenado    | run 1 →      0.350 ms | cmps:          174999 | swaps:              0
  ...

💾 Exportando resultados...
  ✓ output/benchmark_results.csv
  ✓ output/benchmark_stats.csv
```

> ⚠️ **Atenção:** algoritmos O(n²) com entradas grandes (≥ 175.000 no pior caso) podem levar **vários minutos** por execução. O tempo total estimado para 175k elementos é de 20–60 minutos dependendo do hardware.

### 4. Arquivos gerados

| Arquivo | Conteúdo |
|---|---|
| `output/benchmark_results.csv` | Uma linha por execução individual (36 linhas para 4 algoritmos × 3 datasets × 3 runs) |
| `output/benchmark_stats.csv` | Estatísticas agregadas: min, max, média, desvio padrão, mediana (12 linhas) |

---

## Geração de Gráficos

Com os CSVs gerados, execute o script Python:

```bash
# Usando venv (recomendado)
python -m venv .venv
.venv\Scripts\activate        # Windows
source .venv/bin/activate     # Linux/macOS

pip install pandas matplotlib seaborn dataframe-image

python scripts/plot_results.py
```

Os PNGs são salvos em `output/plots_python/<tamanho>/`:

| Arquivo | Descrição |
|---|---|
| `py_tempo_medio<N>.png` | Gráfico de barras agrupado com escala logarítmica (base 10) |
| `py_tabela_resultados<N>.png` | Tabela de estatísticas: min, max, média, desvio padrão, mediana |
| `py_tabela_completa<N>.png` | Tabela completa com todas as execuções individuais |

**Paleta de cores utilizada** (consistente em todos os tamanhos):

| Cor | Dataset |
|---|---|
| 🔵 Azul `#4477AA` | `aleatorio` |
| 🟢 Verde `#228833` | `invertido` |
| 🌸 Rosa `#CC6677` | `ordenado` |

---

## Métricas Coletadas

Cada execução registra os seguintes campos, exportados para CSV:

| Campo | Tipo | Descrição |
|---|---|---|
| `Algorithm` | string | Nome do algoritmo (`BubbleSort`, `SelectionSort`, etc.) |
| `DataType` | string | Tipo do dataset: `ordenado`, `invertido`, `aleatorio` |
| `Run` | int | Índice da execução (1, 2 ou 3) |
| `InputSize` | int | Número de elementos no dataset |
| `DurationNs` | int64 | Duração da ordenação em nanossegundos |
| `DurationMs` | float64 | Duração em milissegundos (3 casas decimais) |
| `Comparisons` | int64 | Total de comparações elemento-a-elemento |
| `Swaps` | int64 | Total de trocas (movimentações de elementos) |
| `MemAllocKB` | float64 | Alocações de heap durante a execução (KB) |

---

## Resultados Parciais

### Tempo médio de execução (escala log₁₀)

As colunas representam: 🔵 aleatorio · 🟢 invertido · 🌸 ordenado

| 175.000 elementos | 250.000 elementos |
|:---:|:---:|
| ![175k](output/plots_python/175k/py_tempo_medio-175.png) | ![250k](output/plots_python/250k/py_tempo_medio250.png) |

| 500.000 elementos | 750.000 elementos |
|:---:|:---:|
| ![500k](output/plots_python/500k/py_tempo_medio500.png) | ![750k](output/plots_python/750k/py_tempo_medio750.png) |

| 1.500.000 elementos |
|:---:|
| ![1500k](output/plots_python/1500k/py_tempo_medio1500.png) |

### Observações experimentais

| Observação | Algoritmo(s) | Explicação |
|---|---|---|
| Tempo ≈ 0 ms para dataset `ordenado` | BubbleSort, InsertionSort | Melhor caso O(n): apenas 1 passagem sem trocas |
| Tempo uniforme independente do dataset | SelectionSort | Sempre faz n(n-1)/2 comparações, independente da entrada |
| Crescimento ~4× ao dobrar n | BubbleSort, SelectionSort, InsertionSort | Confirma empiricamente O(n²) |
| Crescimento ~2,2× ao dobrar n | ShellSort | Consistente com O(n log²n) |
| Shell Sort até 8.000× mais rápido | Shell vs Bubble (1,5M, aleatório) | Diferença de classe de complexidade |
| SelectionSort `ordenado` ≈ demais datasets | SelectionSort | Não há melhor caso — faz todas as comparações sempre |

---

## Adicionando Novos Algoritmos

O projeto foi projetado para receber novos algoritmos **sem modificar nenhum arquivo existente**, exceto `main.go`:

### Passo 1 — Implementar o algoritmo em `benchmark/algorithms.go`

```go
// MergeSorter implementa SortAlgorithm.
type MergeSorter struct{}

func (MergeSorter) Name() string { return "MergeSort" }

func (MergeSorter) Sort(arr []int) SortResult {
    var cmp, swaps int64
    // ... implementação com contagem de cmp e swaps
    return SortResult{arr, cmp, swaps}
}
```

### Passo 2 — Registrar em `cmd/benchmark/main.go`

```go
algorithms := []benchmark.SortAlgorithm{
    benchmark.BubbleSorter{},
    benchmark.SelectionSorter{},
    benchmark.InsertionSorter{},
    benchmark.ShellSorter{},
    benchmark.MergeSorter{},   // ← adicionar aqui
}
```

Nenhum outro arquivo precisa ser alterado. O runner, o exporter e o script Python processam qualquer algoritmo automaticamente.

---

## Próximas Etapas

- [ ] Implementar algoritmos de complexidade O(n log n): Merge Sort, Quick Sort, Heap Sort
- [ ] Implementar algoritmos não-comparativos: Counting Sort, Radix Sort
- [ ] Análise estatística formal (testes de hipótese para comparar médias entre runs)
- [ ] Geração de gráficos de crescimento temporal (tempo × n) para ajuste de curva
- [ ] Redação do artigo científico com todos os resultados consolidados

---

## Instituição

**Universidade Tecnológica Federal do Paraná — UTFPR**  
Disciplina: Pesquisa e Ordenação  
Curso: Ciência da Computação / Engenharia de Computação

---

## Licença

MIT
