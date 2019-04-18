package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:          "merge [flags]",
	SilenceUsage: true,
	Short:        "Merge two TBLNs",
	Long: `Merge two TBLNs.

	The two TBLNs should basically have the same structure.`,
	RunE: merge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	mergeCmd.PersistentFlags().StringP("self-file", "", "", "TBLN self file")
	mergeCmd.PersistentFlags().StringP("self-db", "", "", "Self database driver name")
	mergeCmd.PersistentFlags().StringP("self-dsn", "", "", "Self dsn name")
	mergeCmd.PersistentFlags().StringP("self-schema", "", "", "Self schema Name")
	mergeCmd.PersistentFlags().StringP("self-table", "", "", "Self table name")
	mergeCmd.PersistentFlags().StringP("other-file", "", "", "TBLN other file")
	mergeCmd.PersistentFlags().StringP("other-db", "", "", "other database driver name")
	mergeCmd.PersistentFlags().StringP("other-dsn", "", "", "other dsn name")
	mergeCmd.PersistentFlags().StringP("other-schema", "", "", "other schema Name")
	mergeCmd.PersistentFlags().StringP("other-table", "", "", "other table name")
	mergeCmd.PersistentFlags().BoolP("ignore", "", false, "Ignore when conflict occurs")
	mergeCmd.PersistentFlags().BoolP("update", "", true, "Update when a conflict occurs")
	mergeCmd.PersistentFlags().BoolP("delete", "", false, "Execute update and delete rows not in other")
	mergeCmd.PersistentFlags().BoolP("no-import", "", false, "output without importing into database")
	mergeCmd.PersistentFlags().StringP("output", "o", "", "write to file instead of stdout")
	mergeCmd.PersistentFlags().StringSliceP("hash", "a", []string{"sha256"}, "hash algorithm(sha256 or sha512)")
	mergeCmd.PersistentFlags().StringSliceP("enable-target", "", []string{"name", "type"}, "hash target extra item (all or each name)")
	mergeCmd.PersistentFlags().StringSliceP("disable-target", "", nil, "hash extra items not to be targeted (all or each name)")
	mergeCmd.PersistentFlags().BoolP("sign", "", false, "sign TBLN file")
	mergeCmd.PersistentFlags().BoolP("quiet", "q", false, "do not prompt for password.")
	mergeCmd.PersistentFlags().SortFlags = false
	mergeCmd.Flags().SortFlags = false
}

func getOtherReader(cmd *cobra.Command, args []string) (tbln.Reader, error) {
	var err error
	var otherReader tbln.Reader

	var otherFileName string
	if otherFileName, err = cmd.PersistentFlags().GetString("other-file"); err != nil {
		return nil, err
	}
	if otherFileName == "" && len(args) >= 2 {
		otherFileName = args[0]
	}
	if otherFileName != "" {
		f, err := os.Open(otherFileName)
		if err != nil {
			return nil, err
		}
		otherReader = tbln.NewReader(f)
		return otherReader, nil
	}

	var otherDb, otherDsn string
	if otherDb, err = cmd.PersistentFlags().GetString("other-db"); err != nil {
		return nil, err
	}
	if otherDsn, err = cmd.PersistentFlags().GetString("other-dsn"); err != nil {
		return nil, err
	}

	if otherDb == "" {
		cmd.SilenceUsage = false
		return nil, fmt.Errorf("not enough arguments")
	}
	otherConn, err := db.Open(otherDb, otherDsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", otherDb, err)
	}
	err = otherConn.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", otherDb, err)
	}
	var schema string
	if schema, err = cmd.PersistentFlags().GetString("other-schema"); err != nil {
		return nil, err
	}
	var tableName string
	if tableName, err = cmd.PersistentFlags().GetString("other-table"); err != nil {
		return nil, err
	}
	otherReader, err = otherConn.ReadTable(schema, tableName, nil)
	if err != nil {
		return nil, err
	}
	return otherReader, nil
}

func getSelfReader(cmd *cobra.Command, args []string) (tbln.Reader, error) {
	var err error
	var selfReader tbln.Reader

	var toFileName string
	if toFileName, err = cmd.PersistentFlags().GetString("self-file"); err != nil {
		return nil, err
	}
	if toFileName == "" && len(args) >= 2 {
		toFileName = args[1]
	}
	if toFileName != "" {
		self, err := os.Open(toFileName)
		if err != nil {
			return nil, err
		}
		selfReader = tbln.NewReader(self)
		return selfReader, nil
	}

	var selfDb, selfDsn string
	if selfDb, err = cmd.PersistentFlags().GetString("self-db"); err != nil {
		return nil, err
	}
	if selfDsn, err = cmd.PersistentFlags().GetString("self-dsn"); err != nil {
		return nil, err
	}

	if selfDb == "" {
		cmd.SilenceUsage = false
		return nil, fmt.Errorf("not enough arguments")
	}
	toConn, err := db.Open(selfDb, selfDsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", selfDb, err)
	}
	err = toConn.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", selfDb, err)
	}
	var schema string
	if schema, err = cmd.PersistentFlags().GetString("self-schema"); err != nil {
		return nil, err
	}
	var tableName string
	if tableName, err = cmd.PersistentFlags().GetString("self-table"); err != nil {
		return nil, err
	}
	selfReader, err = toConn.ReadTable(schema, tableName, nil)
	if err != nil {
		return nil, err
	}
	return selfReader, nil
}

func mergeWriteTable(otherReader tbln.Reader, cmd *cobra.Command) error {
	var err error
	var selfDb, selfDsn string
	if selfDb, err = cmd.PersistentFlags().GetString("self-db"); err != nil {
		return err
	}
	if selfDsn, err = cmd.PersistentFlags().GetString("self-dsn"); err != nil {
		return err
	}
	Conn, err := db.Open(selfDb, selfDsn)
	if err != nil {
		return fmt.Errorf("%s: %s", selfDb, err)
	}
	err = Conn.Begin()
	if err != nil {
		return fmt.Errorf("%s: %s", selfDb, err)
	}
	var schema string
	if schema, err = cmd.PersistentFlags().GetString("self-schema"); err != nil {
		return err
	}
	var tableName string
	if tableName, err = cmd.PersistentFlags().GetString("self-table"); err != nil {
		return err
	}

	var ignore, delete bool
	if ignore, err = cmd.PersistentFlags().GetBool("ignore"); err != nil {
		return err
	}
	if delete, err = cmd.PersistentFlags().GetBool("delete"); err != nil {
		return err
	}
	if !ignore {
		err = Conn.MergeTable(schema, tableName, otherReader, delete)
		if err == nil {
			return Conn.Commit()
		}
		log.Printf("Table Not found. Create Table [%s]\n", tableName)
	}
	err = Conn.WriteReader(schema, tableName, otherReader, db.Create, db.OrIgnore)
	if err != nil {
		return err
	}
	return Conn.Commit()
}

func merge(cmd *cobra.Command, args []string) error {
	otherReader, err := getOtherReader(cmd, args)
	if err != nil {
		return err
	}
	var signF bool
	if signF, err = cmd.PersistentFlags().GetBool("sign"); err != nil {
		return err
	}
	if signF && SecFile == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("requires secret key file")
	}

	var noImport bool
	if noImport, err = cmd.PersistentFlags().GetBool("no-import"); err != nil {
		return err
	}
	var selfDb string
	if selfDb, err = cmd.PersistentFlags().GetString("self-db"); err != nil {
		return err
	}
	if selfDb != "" && !noImport {
		return mergeWriteTable(otherReader, cmd)
	}
	mode := tbln.MergeUpdate
	var ignore, delete bool
	if ignore, err = cmd.PersistentFlags().GetBool("ignore"); err != nil {
		return err
	}
	if ignore {
		mode = tbln.MergeIgnore
	}
	if delete, err = cmd.PersistentFlags().GetBool("delete"); err != nil {
		return err
	}
	if delete {
		mode = tbln.MergeDelete
	}
	selfReader, err := getSelfReader(cmd, args)
	if err != nil {
		return err
	}
	if otherReader == nil || selfReader == nil {
		return fmt.Errorf("requires other and self")
	}
	var tb *tbln.Tbln
	tb, err = tbln.MergeAll(selfReader, otherReader, mode)
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
