package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	"github.com/spf13/cobra"
	"github.com/xo/dburl"
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
	mergeCmd.PersistentFlags().StringP("file", "f", "", "self TBLN file")
	mergeCmd.PersistentFlags().StringP("dburl", "d", "", "self database url")
	mergeCmd.PersistentFlags().StringP("schema", "n", "", "self schema Name")
	mergeCmd.PersistentFlags().StringP("table", "t", "", "self table name")
	mergeCmd.PersistentFlags().StringP("other-file", "", "", "other TBLN file")
	mergeCmd.PersistentFlags().StringP("other-dburl", "", "", "other database url")
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

	var url string
	if url, err = cmd.PersistentFlags().GetString("other-dburl"); err != nil {
		return nil, err
	}
	if url == "" {
		return nil, nil
	}
	u, err := dburl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("other:%s: %w", url, err)
	}
	otherConn, err := db.Open(u.Driver, u.DSN)
	if err != nil {
		return nil, fmt.Errorf("other:%s: %w", u.Driver, err)
	}
	err = otherConn.Begin()
	if err != nil {
		return nil, fmt.Errorf("other:%s: %w", u.Driver, err)
	}
	var schema string
	if schema, err = cmd.PersistentFlags().GetString("other-schema"); err != nil {
		return nil, err
	}
	var tableName string
	if tableName, err = cmd.PersistentFlags().GetString("other-table"); err != nil {
		return nil, err
	}
	if tableName == "" {
		if tableName, err = cmd.PersistentFlags().GetString("table"); err != nil {
			return nil, err
		}
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
	if toFileName, err = cmd.PersistentFlags().GetString("file"); err != nil {
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

	var url string
	if url, err = cmd.PersistentFlags().GetString("dburl"); err != nil {
		return nil, err
	}
	if url == "" {
		return nil, nil
	}
	u, err := dburl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("self:%s: %w", url, err)
	}
	toConn, err := db.Open(u.Driver, u.DSN)
	if err != nil {
		return nil, fmt.Errorf("self:%s: %w", u.Driver, err)
	}
	err = toConn.Begin()
	if err != nil {
		return nil, fmt.Errorf("self:%s: %w", u.Driver, err)
	}
	var schema string
	if schema, err = cmd.PersistentFlags().GetString("schema"); err != nil {
		return nil, err
	}
	var tableName string
	if tableName, err = cmd.PersistentFlags().GetString("table"); err != nil {
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
	var url string
	if url, err = cmd.PersistentFlags().GetString("dburl"); err != nil {
		return err
	}
	if url == "" {
		cmd.SilenceUsage = false
		return fmt.Errorf("database url not found")
	}
	u, err := dburl.Parse(url)
	if err != nil {
		return fmt.Errorf("self:%s: %w", url, err)
	}
	conn, err := db.Open(u.Driver, u.DSN)
	if err != nil {
		return fmt.Errorf("self:%s: %w", u.Driver, err)
	}
	err = conn.Begin()
	if err != nil {
		return fmt.Errorf("self:%s: %w", u.Driver, err)
	}

	var schema string
	if schema, err = cmd.PersistentFlags().GetString("schema"); err != nil {
		return err
	}
	var tableName string
	if tableName, err = cmd.PersistentFlags().GetString("table"); err != nil {
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
		err = conn.MergeTable(schema, tableName, otherReader, delete)
		if err == nil {
			return conn.Commit()
		}
		log.Printf("Table Not found. Create Table [%s]\n", tableName)
	}
	err = conn.WriteReader(schema, tableName, otherReader, db.Create, db.OrIgnore)
	if err != nil {
		return err
	}
	return conn.Commit()
}

func merge(cmd *cobra.Command, args []string) error {
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

	var noImport bool
	if noImport, err = cmd.PersistentFlags().GetBool("no-import"); err != nil {
		return err
	}
	var url string
	if url, err = cmd.PersistentFlags().GetString("dburl"); err != nil {
		return err
	}
	if url != "" {
		u, err := dburl.Parse(url)
		if err != nil {
			return fmt.Errorf("self:%s: %w", url, err)
		}
		if u.Driver != "" && !noImport {
			return mergeWriteTable(otherReader, cmd)
		}
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
		cmd.SilenceUsage = false
		return fmt.Errorf("requires other and self")
	}
	var tb *tbln.TBLN
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
