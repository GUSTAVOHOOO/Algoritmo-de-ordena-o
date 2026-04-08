import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import os
import dataframe_image as dfi

# Base directory of the script
script_dir = os.path.dirname(os.path.abspath(__file__))
# Data directory
output_dir = os.path.abspath(os.path.join(script_dir, '..', 'output'))

def load_data():
    stats_path = os.path.join(output_dir, 'benchmark_stats.csv')
    results_path = os.path.join(output_dir, 'benchmark_results.csv')
    stats = pd.read_csv(stats_path)
    results = pd.read_csv(results_path)
    return stats, results

def plot_tempo_medio(stats_df, out_dir):
    plt.figure(figsize=(10, 6))
    
    # Configurar estilo
    sns.set_theme(style="whitegrid")
    
    # Criar barplot agrupado
    ax = sns.barplot(
        data=stats_df,
        x="Algorithm",
        y="MeanMs",
        hue="DataType",
        palette=["#4477AA", "#228833", "#CC6677"]
    )
    
    # Configurar escala logarítmica para o eixo Y
    ax.set_yscale('log')
    
    plt.title('Tempo Médio de Execução por Algoritmo', fontsize=14, pad=15)
    plt.xlabel('Algoritmo', fontsize=12)
    plt.ylabel('Tempo Médio (ms) [escala log10]', fontsize=12)
    
    # Reorganizar legenda
    plt.legend(title='Tipo de Dataset', bbox_to_anchor=(1.05, 1), loc='upper left')
    plt.tight_layout()
    
    plt.savefig(os.path.join(out_dir, 'py_tempo_medio.png'), dpi=300)
    plt.close()

def plot_tabela(stats_df, out_dir):
    # Formata a tabela para melhor apresentação (arredondar decimais, etc)
    df_formatted = stats_df.copy()
    if 'MeanMs' in df_formatted.columns:
        df_formatted['MeanMs'] = df_formatted['MeanMs'].round(3)
    if 'MinMs' in df_formatted.columns:
        df_formatted['MinMs'] = df_formatted['MinMs'].round(3)
    if 'MaxMs' in df_formatted.columns:
        df_formatted['MaxMs'] = df_formatted['MaxMs'].round(3)
    if 'StdDevMs' in df_formatted.columns:
        df_formatted['StdDevMs'] = df_formatted['StdDevMs'].round(3)
    if 'MedianMs' in df_formatted.columns:
        df_formatted['MedianMs'] = df_formatted['MedianMs'].round(3)
        
    # Renomeia as colunas para português
    colunas_pt = {
        'Algorithm': 'Algoritmo',
        'DataType': 'Dataset',
        'Runs': 'Execuções',
        'MinMs': 'Min (ms)',
        'MaxMs': 'Max (ms)',
        'MeanMs': 'Média (ms)',
        'StdDevMs': 'Desv. Padrão (ms)',
        'MedianMs': 'Mediana (ms)'
    }
    df_formatted.rename(columns=colunas_pt, inplace=True)
        
    table_path = os.path.join(out_dir, 'py_tabela_resultados.png')
    # Tenta usar o conversor padrão, se falhar pode precisar do parâmetro table_conversion='matplotlib'
    dfi.export(df_formatted, table_path)

def plot_tabela_completa(results_df, out_dir):
    # Formata a tabela para melhor apresentação (arredondar decimais, etc)
    df_formatted = results_df.copy()
    
    if 'DurationMs' in df_formatted.columns:
        df_formatted['DurationMs'] = df_formatted['DurationMs'].round(3)
    if 'MemAllocKB' in df_formatted.columns:
        df_formatted['MemAllocKB'] = df_formatted['MemAllocKB'].round(3)
        
    # Renomeia as colunas para português
    colunas_pt = {
        'Algorithm': 'Algoritmo',
        'DataType': 'Dataset',
        'Run': 'Execução',
        'InputSize': 'Tamanho Entrada',
        'DurationNs': 'Duração (ns)',
        'DurationMs': 'Duração (ms)',
        'Comparisons': 'Comparações',
        'Swaps': 'Trocas',
        'MemAllocKB': 'Memória (KB)'
    }
    df_formatted.rename(columns=colunas_pt, inplace=True)
        
    table_path = os.path.join(out_dir, 'py_tabela_completa.png')
    dfi.export(df_formatted, table_path)


if __name__ == "__main__":
    out_dir_plots = os.path.join(output_dir, 'plots_python')
    os.makedirs(out_dir_plots, exist_ok=True)
    
    print("📊 Lendo dados do benchmark...")
    stats_df, results_df = load_data()
    
    print("🎨 Gerando gráficos Python com seaborn...")
    plot_tempo_medio(stats_df, out_dir_plots)
    
    print("🖼️ Gerando tabela de estatísticas em PNG...")
    plot_tabela(stats_df, out_dir_plots)
    
    print("🖼️ Gerando tabela completa de resultados em PNG...")
    plot_tabela_completa(results_df, out_dir_plots)
    
    print(f"✅ Arquivos salvos em {out_dir_plots}")
