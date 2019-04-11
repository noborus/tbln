package cmd

import (
	"log"
	"os"

	"github.com/noborus/tbln"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge [TBLN file] [TBLN file]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: merge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}

var data1 = `# String example
; TableName: geh1
; primarykey: | id |
; pg_type: | int | varchar(40) | int |
; name: | id | name | age |
; type: | int | text | int |
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 19 |
| 5 | hoge | 99 |
| 9 | noborus | 46 |
| 11 | ghoge | 4 |
| 21 | aghoge | 4 |
`

var data2 = `# String example
; TableName: newgeh1
; primarykey: | id |
; pg_type: | int | varchar(40) | int |
; name: | id | names | age |
; type: | int | text | int |
| 1 | Beob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 19 |
| 4 | noborus | 46 |
| 9 | noborus | 46 |
| 10 | hoge | 46 |
| 11 | ghoge | 5 |
`

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
	m := tbln.NewMerge(os.Stdout)
	c := tbln.NewCompare(m, sr, dr)
	m.Definition, err = m.MergeDefinition(sr, dr)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Header()
	if err != nil {
		log.Fatal(err)
	}

	err = c.Compare()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
