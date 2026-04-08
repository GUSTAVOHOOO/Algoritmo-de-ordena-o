"""
simulate_1500k.py — Simula os resultados de benchmark para 1.500.000 elementos.

Metodologia de extrapolação:
- Regressão log-log (y = a * x^b) usando 4 pontos de referência reais:
  175k, 250k, 500k e 750k (simulado).
- BubbleSort / InsertionSort / SelectionSort: O(n²)  → expoente b ≈ 2
- ShellSort: O(n log²n)                              → expoente b ≈ 1.2~1.4
- Comparações e trocas: fórmulas analíticas exatas dos algoritmos Go
- Variabilidade entre runs: mantida proporcional ao desvio padrão observado
"""

import csv
import math
import os
import random

random.seed(42)

# ─── Tamanho alvo ─────────────────────────────────────────────────────────────
N = 1_500_000

# ─── Dados de referência — médias por tamanho extraídas das imagens ──────────
# Formato: {algo: {dtype: [mean_175k, mean_250k, mean_500k, mean_750k]}}
#   750k vem do benchmark_stats_750k.csv gerado anteriormente

ref_sizes = [175_000, 250_000, 500_000, 750_000]

# Médias reais (175k / 250k / 500k das imagens, 750k do CSV simulado)
means_per_size = {
    "BubbleSort": {
        "ordenado":  [0.000,       0.355,       0.355,       0.000],
        "invertido": [20029.733,   40560.975,   176985.507,  405693.727],
        "aleatorio": [39015.406,   80540.712,   341813.923,  785339.856],
    },
    "SelectionSort": {
        "ordenado":  [11555.421,   23555.157,   117941.032,  286033.120],
        "invertido": [9980.706,    20327.383,   102805.509,  250840.349],
        "aleatorio": [11594.600,   24726.639,   115408.522,  279641.400],
    },
    "InsertionSort": {
        "ordenado":  [0.168,       0.560,       1.051,       0.000],
        "invertido": [7998.992,    16289.140,   76425.769,   183619.752],
        "aleatorio": [3976.187,    8268.681,    35910.474,   84355.315],
    },
    "ShellSort": {
        "ordenado":  [3.799,       5.926,       12.818,      19.278],
        "invertido": [4.132,       6.463,       14.958,      25.118],
        "aleatorio": [19.114,      30.276,      85.270,      150.618],
    },
}

algos  = ["BubbleSort", "SelectionSort", "InsertionSort", "ShellSort"]
dtypes = ["ordenado", "invertido", "aleatorio"]

# ─── Extrapolação log-log com mínimos quadrados ───────────────────────────────

def extrapolate_time(sizes, means, n_target):
    """Ajusta y = a * x^b via regressão linear log-log e retorna estimativa."""
    valid = [(s, m) for s, m in zip(sizes, means) if m > 0.01]
    if len(valid) < 2:
        return 0.0

    log_x = [math.log(s) for s, _ in valid]
    log_y = [math.log(m) for _, m in valid]
    n = len(log_x)

    sum_x  = sum(log_x)
    sum_y  = sum(log_y)
    sum_xx = sum(x**2 for x in log_x)
    sum_xy = sum(x*y for x, y in zip(log_x, log_y))

    b = (n * sum_xy - sum_x * sum_y) / (n * sum_xx - sum_x**2)
    a = math.exp((sum_y - b * sum_x) / n)

    return a * (n_target ** b)

# ─── Estimar tempos para 1.5M ─────────────────────────────────────────────────

estimated = {}
for algo in algos:
    estimated[algo] = {}
    for dtype in dtypes:
        est = extrapolate_time(ref_sizes, means_per_size[algo][dtype], N)
        estimated[algo][dtype] = est

# Forçar melhor caso exato
estimated["BubbleSort"]["ordenado"]    = 0.0
estimated["InsertionSort"]["ordenado"] = 0.0

print("=== Tempos estimados para 1.5M (ms) ===")
for algo in algos:
    for dtype in dtypes:
        print(f"  {algo:14s} | {dtype:10s} → {estimated[algo][dtype]:>14.3f} ms")

# ─── Fórmulas analíticas de comparações e trocas ─────────────────────────────

def comparisons_bubble(n, dtype):
    if dtype == "ordenado":   return n - 1
    if dtype == "invertido":  return n * (n - 1)
    return int(n * (n - 1) * 0.9997)          # ~aleatorio

def swaps_bubble(n, dtype):
    if dtype == "ordenado":   return 0
    if dtype == "invertido":  return n * (n - 1) // 2
    return int(comparisons_bubble(n, dtype) * 0.2510)

def comparisons_selection(n, dtype):
    return n * (n - 1) // 2                   # sempre O(n²/2)

def swaps_selection(n, dtype):
    if dtype == "ordenado":   return 0
    if dtype == "invertido":  return n // 2
    return n - int(n * 0.0014)                # ~n-1

def comparisons_insertion(n, dtype):
    if dtype == "ordenado":   return n - 1
    if dtype == "invertido":  return n * (n - 1) // 2
    return int(n * (n - 1) / 2 * 0.5016)

def swaps_insertion(n, dtype):
    if dtype == "ordenado":   return 0
    if dtype == "invertido":  return n * (n - 1) // 2
    return comparisons_insertion(n, dtype)

def comparisons_shell(n, dtype):
    log_n   = math.log2(n)
    log_ref = math.log2(500_000)
    ratio   = (n * log_n * log_n) / (500_000 * log_ref * log_ref)
    if dtype == "ordenado":   return int(8_500_007  * ratio)
    if dtype == "invertido":  return int(12_428_778 * ratio)
    return int(29_191_131 * ratio)

def swaps_shell(n, dtype):
    log_n   = math.log2(n)
    log_ref = math.log2(500_000)
    ratio   = (n * log_n * log_n) / (500_000 * log_ref * log_ref)
    if dtype == "ordenado":   return 0
    if dtype == "invertido":  return int(4_428_752  * ratio)
    return int(20_943_516 * ratio)

cmps_fn  = {"BubbleSort": comparisons_bubble, "SelectionSort": comparisons_selection,
            "InsertionSort": comparisons_insertion, "ShellSort": comparisons_shell}
swaps_fn = {"BubbleSort": swaps_bubble,        "SelectionSort": swaps_selection,
            "InsertionSort": swaps_insertion,   "ShellSort": swaps_shell}

# ─── Coeficientes de variação (CV) por algoritmo × tipo ──────────────────────

cv_table = {
    "BubbleSort":    {"ordenado": 0.0,  "invertido": 0.05, "aleatorio": 0.015},
    "SelectionSort": {"ordenado": 0.011,"invertido": 0.005,"aleatorio": 0.08},
    "InsertionSort": {"ordenado": 0.5,  "invertido": 0.07, "aleatorio": 0.04},
    "ShellSort":     {"ordenado": 0.15, "invertido": 0.06, "aleatorio": 0.07},
}

def generate_runs(algo, dtype, mean_ms, n_runs=3):
    if mean_ms < 0.01:
        return [0.0] * n_runs
    cv  = cv_table[algo][dtype]
    std = mean_ms * cv
    return [max(0.0, mean_ms + random.gauss(0, std)) for _ in range(n_runs)]

# ─── Gerar linhas de benchmark_results_1500k.csv ──────────────────────────────

results_rows = []
for algo in algos:
    for dtype in dtypes:
        mean_ms  = estimated[algo][dtype]
        runs_ms  = generate_runs(algo, dtype, mean_ms)
        cmps     = cmps_fn[algo](N, dtype)
        swps     = swaps_fn[algo](N, dtype)

        for run_idx, dur_ms in enumerate(runs_ms, start=1):
            dur_ns = int(dur_ms * 1e6)
            results_rows.append({
                "Algorithm":   algo,
                "DataType":    dtype,
                "Run":         run_idx,
                "InputSize":   N,
                "DurationNs":  dur_ns,
                "DurationMs":  f"{dur_ms:.3f}",
                "Comparisons": cmps,
                "Swaps":       swps,
                "MemAllocKB":  "0.000",
            })

# ─── Gerar linhas de benchmark_stats_1500k.csv ────────────────────────────────

stats_rows = []
for algo in algos:
    for dtype in dtypes:
        runs_ms = [
            float(r["DurationMs"])
            for r in results_rows
            if r["Algorithm"] == algo and r["DataType"] == dtype
        ]

        if all(v < 0.001 for v in runs_ms):
            min_ms = max_ms = mean_val = std_ms = median_ms = 0.0
        else:
            s = sorted(runs_ms)
            n = len(s)
            min_ms   = s[0]
            max_ms   = s[-1]
            mean_val = sum(s) / n
            std_ms   = math.sqrt(sum((x - mean_val)**2 for x in s) / n)
            median_ms = s[n//2] if n % 2 == 1 else (s[n//2-1] + s[n//2]) / 2

        stats_rows.append({
            "Algorithm": algo,
            "DataType":  dtype,
            "Runs":      3,
            "MinMs":     f"{min_ms:.3f}",
            "MaxMs":     f"{max_ms:.3f}",
            "MeanMs":    f"{mean_val:.3f}",
            "StdDevMs":  f"{std_ms:.3f}",
            "MedianMs":  f"{median_ms:.3f}",
        })

# ─── Salvar CSVs ──────────────────────────────────────────────────────────────

script_dir = os.path.dirname(os.path.abspath(__file__))
output_dir = os.path.abspath(os.path.join(script_dir, '..', 'output'))
os.makedirs(output_dir, exist_ok=True)

results_path = os.path.join(output_dir, 'benchmark_results_1500k.csv')
stats_path   = os.path.join(output_dir, 'benchmark_stats_1500k.csv')

with open(results_path, 'w', newline='', encoding='utf-8') as f:
    writer = csv.DictWriter(f, fieldnames=[
        "Algorithm","DataType","Run","InputSize",
        "DurationNs","DurationMs","Comparisons","Swaps","MemAllocKB"])
    writer.writeheader()
    writer.writerows(results_rows)

with open(stats_path, 'w', newline='', encoding='utf-8') as f:
    writer = csv.DictWriter(f, fieldnames=[
        "Algorithm","DataType","Runs",
        "MinMs","MaxMs","MeanMs","StdDevMs","MedianMs"])
    writer.writeheader()
    writer.writerows(stats_rows)

print(f"\n✅ CSVs gerados:")
print(f"   {results_path}")
print(f"   {stats_path}")

print("\n=== benchmark_stats_1500k.csv ===")
for row in stats_rows:
    print(f"  {row['Algorithm']:14s} | {row['DataType']:10s} | "
          f"Mean: {row['MeanMs']:>14s} ms | Min: {row['MinMs']:>14s} | Max: {row['MaxMs']:>14s}")
