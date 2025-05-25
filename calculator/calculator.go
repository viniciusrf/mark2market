package calculator

import (
	"math"
	"time"
)

type Resultado struct {
	ValorInicial     float64   `json:"valor_inicial"`
	TaxaAplicada     float64   `json:"taxa_aplicada"`
	DataCompra       time.Time `json:"data_compra"`
	DataAtualCalculo time.Time `json:"data_calculo"`
	ValorAtualVenda  float64   `json:"valor_venda"`
	ValorFinal       float64   `json:"valor_final"`
	VariacaoPercent  float64   `json:"variacao_percentual"`
	ImpostoCobrado   float64   `json:"imposto_valor"`
	TempoRestante    int       `json:"tempo_restante"`
	TaxaOporMin      float64   `json:"taxa_oportunidade_minima"`
}

func CalculateMarkToMarket(valorInicial, valorVenda, taxaAnual float64, dataCompra, dataFim time.Time) Resultado {
	dataCalculo := time.Now()

	mesesPosse := int(math.Floor(dataCalculo.Sub(dataCompra).Hours() / 24 / 30))
	mesesTotal := int(math.Floor(dataFim.Sub(dataCompra).Hours() / 24 / 30))
	tempoRestante := int(math.Floor(dataFim.Sub(dataCalculo).Hours() / 24 / 30))

	taxaMensal := math.Pow(1+taxaAnual/100, 1.0/12) - 1

	imposto := getAliquotaImposto(mesesPosse)

	valorVenda = valorVenda - ((valorVenda - valorInicial) * (imposto / 100))

	valorFinal := valorInicial * math.Pow(1+taxaMensal, float64(mesesTotal-1))
	// Calcular variações
	variacaoPercent := (valorVenda / valorFinal) * 100

	taxa := (math.Pow(valorFinal/valorVenda, 1/(float64(tempoRestante)/12)) - 1)
	taxaOporMin := taxa * 100

	return Resultado{
		ValorInicial:     valorInicial,
		TaxaAplicada:     taxaAnual,
		DataCompra:       dataCompra,
		DataAtualCalculo: dataCalculo,
		ValorAtualVenda:  valorVenda,
		ValorFinal:       valorFinal,
		VariacaoPercent:  variacaoPercent,
		ImpostoCobrado:   imposto,
		TempoRestante:    tempoRestante,
		TaxaOporMin:      taxaOporMin,
	}
}

func getAliquotaImposto(meses int) float64 {
	dias := meses * 30 // Aproximação de 1 mês = 30 dias

	switch {
	case dias <= 180:
		return 22.5
	case dias <= 360:
		return 20.0
	case dias <= 720:
		return 17.5
	default:
		return 15.0
	}
}
