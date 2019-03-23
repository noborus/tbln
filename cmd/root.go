package cmd

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	// KeyPathFolder is the folder name of the default key path.
	KeyPathFolder = ".tbln"
	// DefaultKeyStore is the storage file name of keystore.
	DefaultKeyStore = "keystore.tbln"
)

// global variable from global flags
var (
	KeyPath  string
	SecFile  string
	PubFile  string
	KeyStore string
	KeyName  string
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
	rootCmd.PersistentFlags().StringVarP(&SecFile, "secfile", "", "", "secret Key file")
	rootCmd.PersistentFlags().StringVarP(&PubFile, "pubfile", "", "", "public Key file")
	rootCmd.PersistentFlags().StringVarP(&KeyStore, "keystore", "", "", "keystore file")
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
		KeyPath = filepath.Join(home, KeyPathFolder)
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
	if SecFile == "" {
		SecFile = filepath.Join(KeyPath, kname+".key")
	} else {
		kname = filepath.Base(SecFile[:len(SecFile)-len(filepath.Ext(SecFile))])
	}
	if PubFile == "" {
		PubFile = filepath.Join(KeyPath, kname+".pub")
	} else {
		kname = filepath.Base(PubFile[:len(PubFile)-len(filepath.Ext(PubFile))])
	}
	if KeyName == "" {
		KeyName = kname
	}
	if KeyStore == "" {
		KeyStore = filepath.Join(KeyPath, DefaultKeyStore)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
