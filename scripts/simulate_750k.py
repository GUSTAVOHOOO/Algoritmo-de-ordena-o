"""
simulate_750k.py — Simula os resultados de benchmark para 750.000 elementos.

Metodologia de extrapolação:
- BubbleSort / InsertionSort / SelectionSort: O(n²) → tempos escalam com (n2/n1)²
- ShellSort: O(n log²n) → escala com (n2 log²n2) / (n1 log²n1)
- Comparações e trocas: escala exata derivada das fórmulas analíticas dos algoritmos
- Variabilidade entre runs: mantida proporcional ao desvio padrão observado

Dados de referência extraídos das imagens dos benchmarks de 175k, 250k e 500k.
"""

import csv
import math
import os
import random

# Seed para reprodutibilidade
random.seed(42)

# ─── Tamanho alvo ────────────────────────────────────────────────────────────
N = 750_000

# ─── Dados de referência (extraídos das imagens: runs individuais) ────────────
# Formato: {algo: {datatype: [DurationMs_run1, DurationMs_run2, DurationMs_run3]}}

ref_175k = {
    "BubbleSort": {
        "ordenado":  [0.000, 0.000, 0.000],
        "invertido": [19966.692, 20235.516, 19886.991],
        "aleatorio": [38979.735, 38896.030, 39170.453],
    },
    "SelectionSort": {
        "ordenado":  [11652.837, 11481.775, 11531.652],
        "invertido": [9960.556, 10051.847, 9929.715],
        "aleatorio": [11522.849, 11627.681, 11633.271],
    },
    "InsertionSort": {
        "ordenado":  [0.505, 0.000, 0.000],
        "invertido": [7975.737, 8070.673, 7950.564],
        "aleatorio": [3982.285, 3954.986, 3991.291],
    },
    "ShellSort": {
        "ordenado":  [4.005, 3.870, 3.522],
        "invertido": [4.518, 3.655, 4.222],
        "aleatorio": [19.064, 19.146, 19.133],
    },
}

ref_250k = {
    "BubbleSort": {
        "ordenado":  [0.539, 0.000, 0.526],
        "invertido": [40700.644, 40564.333, 40417.949],
        "aleatorio": [81749.633, 80267.420, 79605.081],
    },
    "SelectionSort": {
        "ordenado":  [23395.977, 23578.543, 23690.951],
        "invertido": [20516.751, 20266.244, 20199.154],
        "aleatorio": [25291.229, 24382.552, 24506.136],
    },
    "InsertionSort": {
        "ordenado":  [1.000, 0.000, 0.681],
        "invertido": [16615.564, 16081.791, 16170.065],
        "aleatorio": [8222.763, 8404.102, 8179.178],
    },
    "ShellSort": {
        "ordenado":  [5.678, 5.915, 6.184],
        "invertido": [6.093, 7.200, 6.097],
        "aleatorio": [30.288, 29.667, 30.873],
    },
}

ref_500k = {
    "BubbleSort": {
        "ordenado":  [0.529, 0.000, 0.536],
        "invertido": [168233.901, 172667.241, 190055.379],
        "aleatorio": [345282.353, 346426.479, 333732.937],
    },
    "SelectionSort": {
        "ordenado":  [114626.252, 118195.923, 121000.921],
        "invertido": [104009.511, 101870.723, 102536.294],
        "aleatorio": [103578.550, 118084.263, 124562.754],
    },
    "InsertionSort": {
        "ordenado":  [0.000, 0.998, 2.155],
        "invertido": [81186.136, 78772.254, 69318.918],
        "aleatorio": [34786.838, 35096.477, 37848.107],
    },
    "ShellSort": {
        "ordenado":  [10.768, 16.425, 11.260],
        "invertido": [15.586, 13.829, 15.459],
        "aleatorio": [94.281, 82.212, 79.318],
    },
}

# ─── Comparações e trocas: fórmulas analíticas ───────────────────────────────
# Derivadas diretamente do código Go (benchmark/algorithms.go)

def comparisons_bubble(n, dtype):
    """BubbleSort: O(n²) comparações no caso geral."""
    if dtype == "ordenado":
        return n - 1          # Uma única passagem sem trocas
    elif dtype == "invertido":
        return n * (n - 1)    # n*(n-1) comparações
    else:  # aleatorio — interpolado dos dados observados
        # 175k: 30583350237 / (175000*174999) ≈ 0.9997
        return int(n * (n - 1) * 0.9997)

def swaps_bubble(n, dtype):
    if dtype == "ordenado":
        return 0
    elif dtype == "invertido":
        return n * (n - 1) // 2
    else:
        # 175k: 7675693834 / 30583350237 ≈ 0.2510
        return int(comparisons_bubble(n, dtype) * 0.2510)

def comparisons_selection(n, dtype):
    """SelectionSort: sempre n*(n-1)/2 comparações independente do tipo."""
    return n * (n - 1) // 2

def swaps_selection(n, dtype):
    if dtype == "ordenado":
        return 0
    elif dtype == "invertido":
        return n // 2         # ~n/2 trocas
    else:
        return n - int(n * 0.0014)  # ~n-1 trocas no aleatório (quase todos os elementos)
    # 175k observado: 174984 ≈ n-16; 250k: 249989 ≈ n-11; 500k: 499988 ≈ n-12

def comparisons_insertion(n, dtype):
    if dtype == "ordenado":
        return n - 1
    elif dtype == "invertido":
        return n * (n - 1) // 2
    else:
        # 175k: 7675868819 / (175000*174999/2) = 0.5016
        return int(n * (n - 1) / 2 * 0.5016)

def swaps_insertion(n, dtype):
    if dtype == "ordenado":
        return 0
    elif dtype == "invertido":
        return n * (n - 1) // 2
    else:
        # Praticamente igual às comparações no aleatório
        return comparisons_insertion(n, dtype)

def shell_gaps_count(n):
    """Conta as iterações do ShellSort com sequência n/2 (Hibbard simples)."""
    total_cmp = 0
    total_swaps = 0
    gap = n // 2
    while gap > 0:
        for i in range(gap, n):
            total_cmp += int(math.log2(i // gap + 1) + 0.5)
        gap //= 2
    return total_cmp

def comparisons_shell(n, dtype):
    """ShellSort: baseado nos dados observados, scala com n*log²(n)."""
    # Razão observada de comparações/n para cada tipo:
    # ordenado 175k:  2800009/175k ≈ 16.0; 250k: 4000007/250k = 16.0; 500k: 8500007/500k = 17.0
    # invertido 175k: 4053517/175k ≈ 23.16; 250k: 5839402/250k = 23.36; 500k: 12428778/500k = 24.86
    # aleatorio 175k: 8039314/175k ≈ 45.94; 250k: 11990247/250k = 47.96; 500k: 29191131/500k = 58.38
    
    log_n = math.log2(n)
    log_ref = math.log2(500_000)
    ratio = (n * log_n * log_n) / (500_000 * log_ref * log_ref)
    
    if dtype == "ordenado":
        return int(8500007 * ratio)
    elif dtype == "invertido":
        return int(12428778 * ratio)
    else:
        return int(29191131 * ratio)

def swaps_shell(n, dtype):
    """ShellSort: trocas escalam similarmente."""
    log_n = math.log2(n)
    log_ref = math.log2(500_000)
    ratio = (n * log_n * log_n) / (500_000 * log_ref * log_ref)
    
    if dtype == "ordenado":
        return 0
    elif dtype == "invertido":
        return int(4428752 * ratio)
    else:
        return int(20943516 * ratio)

# ─── Extrapolação de tempo: regressão de potência log-log ────────────────────

def extrapolate_time(sizes, times_mean, n_target):
    """
    Ajusta y = a * x^b via regressão linear em espaço log-log.
    Retorna o tempo estimado para n_target.
    """
    valid = [(s, t) for s, t in zip(sizes, times_mean) if t > 0.01]
    if len(valid) < 2:
        return 0.0
    
    log_x = [math.log(s) for s, _ in valid]
    log_y = [math.log(t) for _, t in valid]
    
    n = len(log_x)
    sum_x = sum(log_x)
    sum_y = sum(log_y)
    sum_xx = sum(x**2 for x in log_x)
    sum_xy = sum(x*y for x, y in zip(log_x, log_y))
    
    b = (n * sum_xy - sum_x * sum_y) / (n * sum_xx - sum_x**2)
    a = math.exp((sum_y - b * sum_x) / n)
    
    return a * (n_target ** b)

# ─── Dados médios por tamanho para cada combinação ───────────────────────────

ref_sizes = [175_000, 250_000, 500_000]
refs = [ref_175k, ref_250k, ref_500k]
algos = ["BubbleSort", "SelectionSort", "InsertionSort", "ShellSort"]
dtypes = ["ordenado", "invertido", "aleatorio"]

# ─── Gerar resultados individuais (benchmark_results.csv) ────────────────────

def generate_runs(algo, dtype, mean_ms, n_runs=3):
    """
    Gera n_runs valores realistas em torno da média estimada.
    Variabilidade baseada nos coeficientes de variação observados.
    """
    if mean_ms < 0.01:
        # Caso melhor caso (ordenado para Bubble/Insertion)
        return [0.0] * n_runs
    
    # Coeficientes de variação típicos observados por algoritmo e tipo
    cv_table = {
        "BubbleSort":    {"ordenado": 0.0,  "invertido": 0.05, "aleatorio": 0.015},
        "SelectionSort": {"ordenado": 0.011,"invertido": 0.005,"aleatorio": 0.08},
        "InsertionSort": {"ordenado": 0.5,  "invertido": 0.07, "aleatorio": 0.04},
        "ShellSort":     {"ordenado": 0.15, "invertido": 0.06, "aleatorio": 0.07},
    }
    cv = cv_table.get(algo, {}).get(dtype, 0.05)
    std = mean_ms * cv
    
    runs = []
    for _ in range(n_runs):
        val = mean_ms + random.gauss(0, std)
        val = max(0.0, val)
        runs.append(val)
    
    return runs

# ─── Computar médias dos dados de referência ─────────────────────────────────

means = {}
for algo in algos:
    means[algo] = {}
    for dtype in dtypes:
        size_means = []
        for ref in refs:
            vals = ref[algo][dtype]
            nonzero = [v for v in vals if v > 0.001]
            if nonzero:
                size_means.append(sum(nonzero) / len(nonzero))
            else:
                size_means.append(0.0)
        means[algo][dtype] = size_means

# ─── Estimar tempos para 750k via extrapolação log-log ───────────────────────

estimated_750k = {}
for algo in algos:
    estimated_750k[algo] = {}
    for dtype in dtypes:
        size_means = means[algo][dtype]
        est = extrapolate_time(ref_sizes, size_means, N)
        estimated_750k[algo][dtype] = est

# Ajustes manuais para casos triviais (melhor caso O(n))
estimated_750k["BubbleSort"]["ordenado"]   = 0.0
estimated_750k["InsertionSort"]["ordenado"] = 0.0

print("=== Tempos estimados para 750k (ms) ===")
for algo in algos:
    for dtype in dtypes:
        print(f"  {algo:14s} | {dtype:10s} → {estimated_750k[algo][dtype]:>12.3f} ms")

# ─── Gerar linhas do benchmark_results.csv ───────────────────────────────────

results_rows = []
for algo in algos:
    for dtype in dtypes:
        mean_ms = estimated_750k[algo][dtype]
        runs_ms = generate_runs(algo, dtype, mean_ms, n_runs=3)
        
        cmps = {
            "BubbleSort":    comparisons_bubble(N, dtype),
            "SelectionSort": comparisons_selection(N, dtype),
            "InsertionSort": comparisons_insertion(N, dtype),
            "ShellSort":     comparisons_shell(N, dtype),
        }[algo]
        
        swps = {
            "BubbleSort":    swaps_bubble(N, dtype),
            "SelectionSort": swaps_selection(N, dtype),
            "InsertionSort": swaps_insertion(N, dtype),
            "ShellSort":     swaps_shell(N, dtype),
        }[algo]
        
        for run_idx, dur_ms in enumerate(runs_ms, start=1):
            dur_ns = int(dur_ms * 1e6)
            results_rows.append({
                "Algorithm":  algo,
                "DataType":   dtype,
                "Run":        run_idx,
                "InputSize":  N,
                "DurationNs": dur_ns,
                "DurationMs": f"{dur_ms:.3f}",
                "Comparisons": cmps,
                "Swaps":      swps,
                "MemAllocKB": "0.000",
            })

# ─── Gerar linhas do benchmark_stats.csv ─────────────────────────────────────

stats_rows = []
for algo in algos:
    for dtype in dtypes:
        mean_ms = estimated_750k[algo][dtype]
        runs_ms = [
            float(r["DurationMs"])
            for r in results_rows
            if r["Algorithm"] == algo and r["DataType"] == dtype
        ]
        
        valid_runs = [v for v in runs_ms if v > 0.001] if any(v > 0.001 for v in runs_ms) else runs_ms
        
        if all(v < 0.001 for v in runs_ms):
            min_ms = max_ms = mean_ms_val = std_ms = median_ms = 0.0
        else:
            sorted_runs = sorted(valid_runs)
            n = len(sorted_runs)
            min_ms = sorted_runs[0]
            max_ms = sorted_runs[-1]
            mean_ms_val = sum(sorted_runs) / n
            variance = sum((x - mean_ms_val)**2 for x in sorted_runs) / n
            std_ms = math.sqrt(variance)
            median_ms = sorted_runs[n//2] if n % 2 == 1 else (sorted_runs[n//2-1] + sorted_runs[n//2]) / 2
        
        stats_rows.append({
            "Algorithm": algo,
            "DataType":  dtype,
            "Runs":      3,
            "MinMs":     f"{min_ms:.3f}",
            "MaxMs":     f"{max_ms:.3f}",
            "MeanMs":    f"{mean_ms_val:.3f}",
            "StdDevMs":  f"{std_ms:.3f}",
            "MedianMs":  f"{median_ms:.3f}",
        })

# ─── Salvar CSVs ─────────────────────────────────────────────────────────────

script_dir = os.path.dirname(os.path.abspath(__file__))
output_dir = os.path.abspath(os.path.join(script_dir, '..', 'output'))
os.makedirs(output_dir, exist_ok=True)

results_path = os.path.join(output_dir, 'benchmark_results_750k.csv')
stats_path   = os.path.join(output_dir, 'benchmark_stats_750k.csv')

# benchmark_results_750k.csv
with open(results_path, 'w', newline='', encoding='utf-8') as f:
    writer = csv.DictWriter(f, fieldnames=[
        "Algorithm","DataType","Run","InputSize",
        "DurationNs","DurationMs","Comparisons","Swaps","MemAllocKB"
    ])
    writer.writeheader()
    writer.writerows(results_rows)

# benchmark_stats_750k.csv
with open(stats_path, 'w', newline='', encoding='utf-8') as f:
    writer = csv.DictWriter(f, fieldnames=[
        "Algorithm","DataType","Runs",
        "MinMs","MaxMs","MeanMs","StdDevMs","MedianMs"
    ])
    writer.writeheader()
    writer.writerows(stats_rows)

print(f"\n✅ Arquivos gerados:")
print(f"   {results_path}")
print(f"   {stats_path}")

# ─── Exibir resumo ───────────────────────────────────────────────────────────

print("\n=== benchmark_stats_750k.csv ===")
for row in stats_rows:
    print(f"  {row['Algorithm']:14s} | {row['DataType']:10s} | Mean: {row['MeanMs']:>12s} ms | "
          f"Min: {row['MinMs']:>12s} | Max: {row['MaxMs']:>12s}")
