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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root")
	},
}

func init() {
	rootCmd.AddCommand()
}

func Execute() {
	rootCmd.PersistentFlags().StringVar(&dbName, "db", "", "database name")
	rootCmd.PersistentFlags().StringVar(&dsn, "dsn", "", "dsn name")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
