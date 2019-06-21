package cmd

import (
	"fmt"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// exceptCmd represents the except command
var exceptCmd = &cobra.Command{
	Use:          "except",
	SilenceUsage: true,
	Short:        "Except for other TBLN rows from self TBLN",
	Long: `Except for other TBLN rows from self TBLN

	The two TBLNs should basically have the same structure.`,
	RunE: except,
}

func init() {
	rootCmd.AddCommand(exceptCmd)
	exceptCmd.PersistentFlags().StringP("file", "f", "", "self TBLN File")
	exceptCmd.PersistentFlags().StringP("dburl", "d", "", "self database url")
	exceptCmd.PersistentFlags().StringP("schema", "n", "", "self schema Name")
	exceptCmd.PersistentFlags().StringP("table", "t", "", "self table name")
	exceptCmd.PersistentFlags().StringP("other-file", "", "", "other TBLN File")
	exceptCmd.PersistentFlags().StringP("other-dburl", "", "", "other database url")
	exceptCmd.PersistentFlags().StringP("other-schema", "", "", "other schema Name")
	exceptCmd.PersistentFlags().StringP("other-table", "", "", "other table name")
	exceptCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")
	exceptCmd.PersistentFlags().StringSliceP("hash", "a", []string{"sha256"}, "hash algorithm(sha256 or sha512)")
	exceptCmd.PersistentFlags().StringSliceP("enable-target", "", []string{"name", "type"}, "hash target extra item (all or each name)")
	exceptCmd.PersistentFlags().StringSliceP("disable-target", "", nil, "hash extra items not to be targeted (all or each name)")
	exceptCmd.PersistentFlags().BoolP("sign", "", false, "sign TBLN file")
	exceptCmd.PersistentFlags().BoolP("quiet", "q", false, "do not prompt for password.")
	exceptCmd.PersistentFlags().SortFlags = false
	exceptCmd.Flags().SortFlags = false
}

func except(cmd *cobra.Command, args []string) error {
	var err error
	var signF bool
	if signF, err = cmd.PersistentFlags().GetBool("sign"); err != nil {
		return err
	}
	if signF && SecFile == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires secret key file")
	}

	otherReader, err := getOtherReader(cmd, args)
	if err != nil {
		return err
	}
	selfReader, err := getSelfReader(cmd, args)
	if err != nil {
		return err
	}
	if otherReader == nil || selfReader == nil {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires from and self")
	}
	var tb *tbln.TBLN
	tb, err = tbln.ExceptAll(selfReader, otherReader)
	if err != nil {
		return err
	}

	tb, err = hashFile(tb, cmd)
	if err != nil {
		return err
	}
	if signF {
		tb, err = signFile(tb, cmd)
		if err != nil {
			return err
		}
	}

	return outputFile(tb, cmd)
}
