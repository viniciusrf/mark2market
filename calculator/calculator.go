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

	taxaMensal := TaxaAnoEmMesesPercToDec(taxaAnual)

	imposto := GetAliquotaImpostoRF(mesesPosse)

	valorVenda = RendimentoMenosImposto(valorInicial, valorVenda, imposto)

	valorFinal := CalcJurosCompostos(valorInicial, taxaMensal, mesesTotal-1)

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

func TempoAteMeta(valorInicial, valorDesejado, taxaMensal float64, aporteMensal float64) (int, float64) {

	if valorInicial >= valorDesejado {
		return 0, valorInicial
	}

	if taxaMensal <= 0 {
		if aporteMensal <= 0 {
			return -1, valorInicial
		}
		months := math.Ceil((valorDesejado - valorInicial) / aporteMensal)
		return int(months), valorInicial + (aporteMensal * months)
	}

	var months int
	valorAtual := valorInicial

	for valorAtual < valorDesejado {
		months++
		valorAtual *= (1 + taxaMensal)
		valorAtual += aporteMensal

		if months > 1200 {
			return -1, valorAtual
		}
	}
	return months, valorAtual
}

func TaxaAnoEmMesesPercToDec(taxaAnual float64) float64 {

	taxaMensal := math.Pow(1+taxaAnual/100, 1.0/12) - 1

	return taxaMensal
}

func TaxaMesesEmAnoPercToDec(taxaMensalPercentual float64) float64 {

	taxaMensal := taxaMensalPercentual / 100
	taxaAnual := (math.Pow(1+taxaMensal, 12) - 1) * 100
	return taxaAnual
}

func GetAliquotaImpostoRF(meses int) float64 {
	dias := meses * 30

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

func CalcJurosCompostos(initial, taxa float64, time int) float64 {

	valorFinal := initial * math.Pow(1+taxa, float64(time))

	return valorFinal

}

func RendimentoMenosImposto(initial, final, imposto float64) float64 {

	valorfinal := final - ((final - initial) * (imposto / 100))

	return valorfinal

}
