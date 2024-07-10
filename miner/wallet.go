package miner

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ssh/terminal"
	"log"
)

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
