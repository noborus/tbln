package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// global variable from global flags
var (
	keypath string
	seckey  string
	pubfile string
	keyname string
)

var rootCmd = &cobra.Command{
	Use:   "tbln",
	Short: "tbln is command-line TBLN processor",
	Long: `Import/Export TBLN file and RDBMS table.
Also sign and verify the TBLN file.`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&keypath, "keypath", "", "", "key store path")
	rootCmd.PersistentFlags().StringVarP(&seckey, "seckey", "", "", "Secret Key File")
	rootCmd.PersistentFlags().StringVarP(&pubfile, "pubfile", "", "", "public Key File")
	rootCmd.PersistentFlags().StringVarP(&keyname, "keyname", "k", "", "key name")
	rootCmd.AddCommand()
}

func initConfig() {
	keyPathStr := os.Getenv("TBLN_KEYPATH")
	if keyPathStr != "" {
		keypath = keyPathStr
	}
	if keypath == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		keypath = filepath.Join(home, `.tbln`)
	}
	kname := ""
	if keyname != "" {
		kname = keyname
	} else {
		user, err := user.Current()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		kname = user.Username
	}
	if seckey == "" {
		seckey = filepath.Join(keypath, kname+".key")
	} else {
		kname = filepath.Base(seckey[:len(seckey)-len(filepath.Ext(seckey))])
	}
	if pubfile == "" {
		pubfile = filepath.Join(keypath, kname+".pub")
	} else {
		kname = filepath.Base(pubfile[:len(pubfile)-len(filepath.Ext(pubfile))])
	}
	if keyname == "" {
		keyname = kname
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
