package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/noborus/tbln/cmd/key"
	"github.com/noborus/tbln/cmd/keystore"

	"github.com/spf13/cobra"
)

// genkeyCmd represents the generate key pair command
var genkeyCmd = &cobra.Command{
	Use:          "genkey [Key Name]",
	SilenceUsage: true,
	Short:        "Generate a new key pair",
	Long: `Generate a new Ed25519 keypair.
Write private key(Key Name+".key") and public key(Key Name+".pub") files
based on the Key Name.`,
	RunE: genkey,
}

func init() {
	genkeyCmd.PersistentFlags().BoolP("force", "", false, "Overwrites even if the file already exists.")

	rootCmd.AddCommand(genkeyCmd)
}

func genkey(cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		KeyName = args[0]
	}
	var overwrite bool
	if overwrite, err = cmd.PersistentFlags().GetBool("force"); err != nil {
		return err
	}
	err = generateKey(KeyName, overwrite)
	if err != nil {
		return err
	}
	return nil
}

func generateKey(keyName string, overwrite bool) error {
	if _, err := os.Stat(KeyPath); os.IsNotExist(err) {
		log.Printf("mkdir [%s]", KeyPath)
		err := os.Mkdir(KeyPath, 0700)
		if err != nil {
			return err
		}
	}
	_, err := os.Stat(SecFile)
	if !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("%s file already exists", SecFile)
	}
	_, err = os.Stat(PubFile)
	if !os.IsNotExist(err) && !overwrite {
		return fmt.Errorf("%s file already exists", PubFile)
	}

	public, private, err := key.GenerateKey()
	if err != nil {
		return err
	}

	err = key.WritePrivateFile(SecFile, keyName, private)
	if err != nil {
		return fmt.Errorf("%s: %w", SecFile, err)
	}
	log.Printf("write %s file\n", SecFile)

	err = key.WritePublicFile(PubFile, keyName, public)
	if err != nil {
		return fmt.Errorf("generate key create: %s: %w", PubFile, err)
	}
	log.Printf("write %s file\n", PubFile)

	err = keystore.Register(KeyStore, keyName, public)
	if err != nil {
		return err
	}
	log.Printf("register %s file\n", KeyStore)
	return nil
}
