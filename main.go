package main

import (
	"os"

	"github.com/spf13/cobra"
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
	rootCmd.AddCommand(jurosCompostos())
	rootCmd.AddCommand(tempoAte())
	rootCmd.AddCommand(devoParcelar())

	return rootCmd
}
