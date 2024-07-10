package commands

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

// GenNodeKeyCmd allows the generation of a node key. It prints node's ID to
// the standard output.
var ServiceStart = &cobra.Command{
	Use:     "service-start",
	Aliases: []string{"service_start"},
	Short:   "Generate task commitment to the network",
	PreRun:  deprecateSnakeCase,
	RunE:    serviceStartWithConfig,
}

func serviceStartWithConfig(cmd *cobra.Command, args []string) error {
	const keyfilePath = "path/to/your/keystore"

	// Check if the key file exists
	_, err := os.Stat(keyfilePath)
	if os.IsNotExist(err) {
		fmt.Println("No keyfile found. Creating a new key...")
		password := promptPassword(true) // Ask for password with confirmation
		createNewKey(keyfilePath, password)
	} else if err != nil {
		return fmt.Errorf("Failed to check key file: %v", err)
	} else {
		fmt.Println("Keyfile found. Loading key...")
		password := promptPassword(false) // Ask for password without confirmation
		loadKey(keyfilePath, password)
	}

	return nil
}

func promptPassword(confirm bool) string {
	fmt.Print("Enter Password: ")
	passwordBytes, err := terminal.ReadPassword(0)
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	fmt.Println()

	if confirm {
		fmt.Print("Confirm Password: ")
		confirmPasswordBytes, err := terminal.ReadPassword(0)
		if err != nil {
			log.Fatalf("Failed to read password confirmation: %v", err)
		}
		fmt.Println()

		if string(passwordBytes) != string(confirmPasswordBytes) {
			log.Fatal("Passwords do not match.")
		}
	}

	return string(passwordBytes)
}

func createNewKey(filePath, password string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	ks := keystore.NewKeyStore(filePath, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.ImportECDSA(privateKey, password)
	if err != nil {
		log.Fatalf("Failed to save the key file: %v", err)
	}

	fmt.Printf("New key created: %s\n", account.Address.Hex())
}

func loadKey(filePath, password string) {
	ks := keystore.NewKeyStore(filePath, keystore.StandardScryptN, keystore.StandardScryptP)
	if len(ks.Accounts()) == 0 {
		log.Fatal("No accounts found in the key store.")
	}
	account := ks.Accounts()[0]

	err := ks.Unlock(account, password)
	if err != nil {
		log.Fatalf("Failed to unlock the account: %v", err)
	}

	fmt.Printf("Account %s loaded\n", account.Address.Hex())
}
