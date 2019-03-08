package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tbln",
	Short: "tbln is command-line TBLN processor",
	Long: `Import/Export TBLN file and RDBMS table.
Also sign and verify the TBLN file.`,
}

func init() {
	rootCmd.AddCommand()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
