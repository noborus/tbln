package cmd

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// global variable from global flags
var (
	KeyPath string
	SecKey  string
	PubFile string
	KeyName string

	Verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "tbln",
	Short: "tbln is command-line TBLN processor",
	Long: `Import/Export TBLN file and RDBMS table.
Supports digital signatures and verification for TBLN files.`,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&KeyPath, "keypath", "", "", "key store path")
	rootCmd.PersistentFlags().StringVarP(&SecKey, "seckey", "", "", "secret Key File")
	rootCmd.PersistentFlags().StringVarP(&PubFile, "PubFile", "", "", "public Key File")
	rootCmd.PersistentFlags().StringVarP(&KeyName, "keyname", "k", "", "key name")
	rootCmd.AddCommand()
}

func initConfig() {
	KeyPathStr := os.Getenv("TBLN_KEYPATH")
	if KeyPathStr != "" {
		KeyPath = KeyPathStr
	}
	if KeyPath == "" {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		KeyPath = filepath.Join(home, `.tbln`)
	}
	kname := ""
	if KeyName != "" {
		kname = KeyName
	} else {
		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		kname = user.Username
	}
	if SecKey == "" {
		SecKey = filepath.Join(KeyPath, kname+".key")
	} else {
		kname = filepath.Base(SecKey[:len(SecKey)-len(filepath.Ext(SecKey))])
	}
	if PubFile == "" {
		PubFile = filepath.Join(KeyPath, kname+".pub")
	} else {
		kname = filepath.Base(PubFile[:len(PubFile)-len(filepath.Ext(PubFile))])
	}
	if KeyName == "" {
		KeyName = kname
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
