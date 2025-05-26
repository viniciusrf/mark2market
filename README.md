# Mark2Market - Ferramenta de Cálculos Financeiros

## Visão Geral

O Mark2Market é uma ferramenta CLI (Command Line Interface) para cálculos financeiros avançados, incluindo:

- Cálculos de mark-to-market para títulos de investimento
- Projeções com juros compostos
- Estimativas de tempo para atingir metas financeiras

## Instalação

1. Certifique-se de ter o Go instalado (versão 1.16 ou superior)
2. Clone o repositório:
   ```bash
   git clone https://github.com/viniciusrf/mark2market.git
   cd mark2market

## Comandos

1. mark2market - Calcula os valores basicos para avaliação de venda de títulos de renda fixa antecipadamente
    ``` bash
    ./mark2market mark2market -c <valor_inicial> -a <valor_atual> -r <taxa_anual%> -i <data_inicio-MM/AAAA> -f <data_fim-MM/AAAA>

2. jusrosCompostos - Calcula o valor final de um investimento, dado um valor inicial e uma taxa de juros
    ``` bash
    ./mark2market jurosCompostos -c <valor_inicial> -r <taxa_anual%> -f <data_fim-DD/MM/AAAA>

3. tempoAte - Calcula o tempo mínimo para atingimento de uma meta, a partir de um valor inicial, taxa de juros média e valor de aportes mensaid
    ``` bash
    ./mark2market tempoAte -c <valor_inicial> -f <valor_desejado> -r <taxa_anual%> -a <aportes_mensais>