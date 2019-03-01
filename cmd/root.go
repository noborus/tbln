package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// global variable from global flags
var (
	dbName string
	dsn    string
)

var rootCmd = &cobra.Command{
	Use:   "tbln",
	Short: "tbln cli tool",
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbName, "db", "", "database name")
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "dsn name")
	rootCmd.AddCommand()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
