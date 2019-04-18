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
	exceptCmd.PersistentFlags().StringP("self-file", "", "", "TBLN self File")
	exceptCmd.PersistentFlags().StringP("self-db", "", "", "Self database driver name")
	exceptCmd.PersistentFlags().StringP("self-dsn", "", "", "Self dsn name")
	exceptCmd.PersistentFlags().StringP("self-schema", "", "", "Self schema Name")
	exceptCmd.PersistentFlags().StringP("self-table", "", "", "Self table name")
	exceptCmd.PersistentFlags().StringP("other-file", "", "", "TBLN other File")
	exceptCmd.PersistentFlags().StringP("other-db", "", "", "Other database driver name")
	exceptCmd.PersistentFlags().StringP("other-dsn", "", "", "Other dsn name")
	exceptCmd.PersistentFlags().StringP("other-schema", "", "", "Other schema Name")
	exceptCmd.PersistentFlags().StringP("other-table", "", "", "Other table name")
	exceptCmd.PersistentFlags().StringP("output", "o", "", "Write to file instead of stdout")
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
		return fmt.Errorf("requires from and self")
	}
	var tb *tbln.Tbln
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