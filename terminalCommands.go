package main

import (
	"fmt"
	"math"
	"time"

	"github.com/spf13/cobra"
	calculator "github.com/viniciusrf/mark2market/calculator"
)

func mark2marketValue() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mark2market",
		Short: "Realiza o cálculo mark-to-market para avaliar se deve vender o título",
		Long:  `Realiza o cálculo de mark-to-market com base no valor inicial e na taxa fornecida e valor atual.`,
		Run:   mark2marketValueHandler,
	}

	cmd.Flags().Float64P("initial", "c", 0, "Valor inicial para cálculo (obrigatório)")
	cmd.Flags().Float64P("atual", "a", 0, "Valor atual de venda para cálculo (obrigatório)")
	cmd.Flags().Float64P("initialRate", "r", 0, "Taxa a adiquirida a.a. no início (em porcentagem, obrigatório)")
	cmd.Flags().StringP("dataInicio", "i", "", "Mes/Ano da compra (MM/AAAA, obrigatório)")
	cmd.Flags().StringP("dataFim", "f", "", "Mes/Ano da compra (MM/AAAA, obrigatório)")

	cmd.MarkFlagRequired("initial")
	cmd.MarkFlagRequired("atual")
	cmd.MarkFlagRequired("dataInicio")
	cmd.MarkFlagRequired("dataFim")
	cmd.MarkFlagRequired("initialRate")

	return cmd
}

func mark2marketValueHandler(cmd *cobra.Command, args []string) {
	valorInicial, _ := cmd.Flags().GetFloat64("initial")
	valorVenda, _ := cmd.Flags().GetFloat64("atual")
	taxa, _ := cmd.Flags().GetFloat64("initialRate")
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
	fmt.Printf("Taxa minima para mesmo ganho: %15.2f%%\n", result.TaxaOporMin)
}

func jurosCompostos() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jurosCompostos",
		Short: "Realiza o cálculo de juros compostos em um tempo",
		Long:  `Realiza o cálculo  de juros compostos baseado no valor inicial e taxa ao ano em %`,
		Run:   jurosCompostosHandler,
	}

	cmd.Flags().Float64P("initial", "c", 0, "Valor inicial para cálculo (obrigatório)")
	cmd.Flags().StringP("dataFim", "f", "", "Dia/Mes/Ano da compra (DD/MM/AAAA, obrigatório)")
	cmd.Flags().Float64P("initialRate", "r", 0, "Taxa a adiquirida a.a. no início (em porcentagem, obrigatório)")

	cmd.MarkFlagRequired("initial")
	cmd.MarkFlagRequired("dataFim")
	cmd.MarkFlagRequired("initialRate")

	return cmd
}

func jurosCompostosHandler(cmd *cobra.Command, args []string) {
	valorInicial, _ := cmd.Flags().GetFloat64("initial")
	taxa, _ := cmd.Flags().GetFloat64("initialRate")
	dataFim, _ := cmd.Flags().GetString("dataFim")

	dataFimConv, err := time.Parse("01/01/2006", dataFim)
	if err != nil {
		fmt.Printf("Erro ao parsear data: %v\n", err)
		return
	}
	mesesTotal := int(math.Floor(time.Until(dataFimConv).Hours() / 24 / 30))
	taxa = calculator.TaxaAnoEmMesesPercToDec(taxa)

	result := calculator.CalcJurosCompostos(valorInicial, taxa, mesesTotal)

	fmt.Printf("\nCalculo de Juros Compostos\n")
	fmt.Printf("-------------------------\n")
	fmt.Printf("Valor inicial: %15.2f\n", valorInicial)
	fmt.Printf("Taxa: %21.2f%%\n", taxa)
	fmt.Printf("Tempo de permanencia: %d meses\n", mesesTotal)
	fmt.Printf("Valor final: %15.2f\n", result)

}

func tempoAte() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tempoAte",
		Short: "Realiza o cálculo do tempo até atingir um valor desejado",
		Long:  `Realiza o cálculo  do tempo até atingir um valor Y, a partir de um valor X e uma taxa média T`,
		Run:   tempoAteHandler,
	}

	cmd.Flags().Float64P("initial", "c", 0, "Valor inicial para cálculo (obrigatório)")
	cmd.Flags().Float64P("final", "f", 0, "Valor desejado para cálculo (obrigatório)")
	cmd.Flags().Float64P("initialRate", "r", 0, "Taxa a média esperada a.a. (em porcentagem, obrigatório)")
	cmd.Flags().Float64P("aportes", "a", 0, "Valor aportado mensalmente (obrigatório)")

	cmd.MarkFlagRequired("initial")
	cmd.MarkFlagRequired("final")
	cmd.MarkFlagRequired("initialRate")
	cmd.MarkFlagRequired("aportes")

	return cmd
}

func tempoAteHandler(cmd *cobra.Command, args []string) {
	valorInicial, _ := cmd.Flags().GetFloat64("initial")
	valorDesejado, _ := cmd.Flags().GetFloat64("final")
	taxa, _ := cmd.Flags().GetFloat64("initialRate")
	aportesMensais, _ := cmd.Flags().GetFloat64("aportes")

	taxaMensal := calculator.TaxaAnoEmMesesPercToDec(taxa)

	meses, valorFinal := calculator.TempoAteMeta(
		valorInicial,
		valorDesejado,
		taxaMensal,
		aportesMensais,
	)

	if meses < 0 {
		fmt.Println("Não será possível atingir o valor alvo com os parâmetros fornecidos")
	} else {
		anos := meses / 12
		meses = meses % 12

		fmt.Printf("\nCalculo de tempo até a meta\n")
		fmt.Printf("-------------------------\n")
		fmt.Printf("Valor inicial: %15.2f\n", valorInicial)
		fmt.Printf("Valor desejado: %15.2f\n", valorDesejado)
		fmt.Printf("Taxa: %21.2f%%\n\n", taxa)
		fmt.Printf("Tempo mínimo até a meta: %d anos e %d meses\n", anos, meses)
		fmt.Printf("Valor final: %15.2f\n", valorFinal)
	}

}

func devoParcelar() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devoParcelar",
		Short: "Realiza o cálculo comparativo entre parcelamento e à vista",
		Long:  `Realiza o cálculo  comparativo entre parcelamento e compra à vista para analisar a melhor opção`,
		Run:   devoParcelarHandler,
	}

	cmd.Flags().Float64P("vista", "v", 0, "Valor à vista para cálculo (obrigatório)")
	cmd.Flags().Float64P("prazo", "p", 0, "Valor a prazo para cálculo (obrigatório)")
	cmd.Flags().IntP("parc", "m", 0, "Quantidade de parcelas")
	cmd.Flags().Float64P("taxa", "r", 0, "Taxa de juros base para calculo (em porcentagem, obrigatório)")

	cmd.MarkFlagRequired("vista")
	cmd.MarkFlagRequired("prazo")
	cmd.MarkFlagRequired("parc")
	cmd.MarkFlagRequired("taxa")

	return cmd
}

func devoParcelarHandler(cmd *cobra.Command, args []string) {
	aVista, _ := cmd.Flags().GetFloat64("vista")
	aPrazo, _ := cmd.Flags().GetFloat64("prazo")
	taxa, _ := cmd.Flags().GetFloat64("taxa")
	meses, _ := cmd.Flags().GetInt("parc")

	taxa = calculator.TaxaAnoEmMesesPercToDec(taxa)

	value, message := calculator.DeveParcelar(aVista, aPrazo, taxa, meses)

	fmt.Printf("\nCalculo de Juros Compostos\n")
	fmt.Printf("-------------------------\n")
	fmt.Printf("Valor à Vista: %15.2f\n", aVista)
	fmt.Printf("Valor à Prazo: %15.2f\n", aPrazo)
	fmt.Printf("Taxa: %21.2f%%\n", taxa)
	fmt.Printf("Tempo de permanencia: %d meses\n", meses)
	fmt.Printf("Valor economizado: %15.2f\n", value)
	fmt.Printf("Deve parcelar: %s\n", message)

}
