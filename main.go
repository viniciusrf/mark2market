package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	calculator "github.com/viniciusrf/mark2market/calculator"
)

func main() {
	rootCmd := NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "mark2market",
		Short: "Ferramenta para cálculo mark-to-market",
		Long: `Uma ferramenta de linha de comando para realizar cálculos
financeiros de mark-to-market (avaliação a mercado).`,
	}

	rootCmd.AddCommand(mark2marketValue())

	return rootCmd
}

func mark2marketValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shouldSell",
		Short: "Realiza o cálculo mark-to-market para avaliar se deve vender o título",
		Long:  `Realiza o cálculo de mark-to-market com base no valor inicial e na taxa fornecida e valor atual.`,
		Run:   mark2marketValueHandler,
	}

	cmd.Flags().Float64P("initial", "c", 0, "Valor inicial para cálculo (obrigatório)")
	cmd.Flags().Float64P("atual", "a", 0, "Valor atual de venda para cálculo (obrigatório)")
	cmd.Flags().Float64P("rate", "r", 0, "Taxa a adiquirida a.a. no início (em porcentagem, obrigatório)")
	cmd.Flags().StringP("dataInicio", "i", "", "Mes/Ano da compra (MM/AAAA, obrigatório)")
	cmd.Flags().StringP("dataFim", "f", "", "Mes/Ano da compra (MM/AAAA, obrigatório)")
	cmd.MarkFlagRequired("initial")
	cmd.MarkFlagRequired("atual")
	cmd.MarkFlagRequired("dataInicio")
	cmd.MarkFlagRequired("initialRate")

	return cmd
}

func mark2marketValueHandler(cmd *cobra.Command, args []string) {
	valorInicial, _ := cmd.Flags().GetFloat64("initial")
	valorVenda, _ := cmd.Flags().GetFloat64("atual")
	taxa, _ := cmd.Flags().GetFloat64("rate")
	dataInicial, _ := cmd.Flags().GetString("dataInicio")
	dataFim, _ := cmd.Flags().GetString("dataFim")

	dataInicialConv, err := time.Parse("01/2006", dataInicial)
	if err != nil {
		fmt.Printf("Erro ao parsear data: %v\n", err)
		return
	}

	dataFimConv, err := time.Parse("01/2006", dataFim)
	if err != nil {
		fmt.Printf("Erro ao parsear data: %v\n", err)
		return
	}

	result := calculator.CalculateMarkToMarket(valorInicial, valorVenda, taxa, dataInicialConv, dataFimConv)

	fmt.Printf("\nCalculo de marcação a mercado\n")
	fmt.Printf("-------------------------\n")
	fmt.Printf("Valor inicial: %15.2f\n", result.ValorInicial)
	fmt.Printf("Taxa: %21.2f%%\n", result.TaxaAplicada)
	fmt.Printf("Valor venda (atual): %15.2f\n", valorVenda)
	fmt.Printf("Imposto cobrado: %15.2f\n", result.ImpostoCobrado)
	fmt.Printf("Valor venda (corrigido): %15.2f\n", result.ValorAtualVenda)
	fmt.Printf("Valor Final: %15.2f\n", result.ValorFinal)
	fmt.Printf("Porcentagem do total: %15.2f\n", result.VariacaoPercent)
	fmt.Printf("Tempo Restante: %d meses\n", result.TempoRestante)
	fmt.Printf("Taxa minima para mesmo ganho: %15.2f\n", result.TaxaOporMin)

	// 	ValorInicial:       valorInicial,
	// TaxaAplicada:       taxaAnual,
	// DataCompra:         dataCompra,
	// DataAtualCalculo:   dataCalculo,
	// ValorAtualEsperado: valorAtual,
	// ValorAtualVenda:    valorVenda,
	// VariacaoPercent:    variacaoPercent,
	// ImpostoCobrado:     imposto,
	// TempoRestante:      tempoRestante,
	// TaxaOporMin:        taxaOporMin,
}
