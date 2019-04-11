package cmd

import (
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge [flags] [TBLN file] [TBLN file]",
	Short: "Merge two TBLNs",
	Long: `Merge two TBLNs.

	The two TBLNs should basically have the same structure.`,
	RunE: merge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")
}

func merge(cmd *cobra.Command, args []string) error {
	var err error
	srcFileName := args[0]
	dstFileName := args[1]
	src, err := os.Open(srcFileName)
	if err != nil {
		return err
	}
	dst, err := os.Open(dstFileName)
	if err != nil {
		return err
	}
	sr := tbln.NewReader(src)
	dr := tbln.NewReader(dst)

	at, err := tbln.MergeAll(sr, dr)
	if err != nil {
		return err
	}
	return outputFile(at, cmd)
}
