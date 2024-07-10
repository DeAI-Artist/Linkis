package miner

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/ssh/terminal"
)

// PromptPassword requests a password from the user and optionally confirms it.
func PromptPassword(confirm bool) (string, error) {
	fmt.Print("Enter Password: ")
	passwordBytes, err := terminal.ReadPassword(0)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	fmt.Println()

	if confirm {
		fmt.Print("Confirm Password: ")
		confirmPasswordBytes, err := terminal.ReadPassword(0)
		if err != nil {
			return "", fmt.Errorf("failed to read password confirmation: %v", err)
		}
		fmt.Println()

		if string(passwordBytes) != string(confirmPasswordBytes) {
			return "", fmt.Errorf("passwords do not match")
		}
	}

	return string(passwordBytes), nil
}

// CreateNewKey generates a new key and stores it in a keystore file.
func CreateNewKey(filePath, password string) error {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	ks := keystore.NewKeyStore(filePath, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.ImportECDSA(privateKey, password)
	if err != nil {
		return fmt.Errorf("failed to save the key file: %v", err)
	}

	fmt.Printf("New key created: %s\n", account.Address.Hex())
	return nil
}

// LoadKey loads a key from a keystore file using the provided password.
func LoadKey(filePath, password string) error {
	ks := keystore.NewKeyStore(filePath, keystore.StandardScryptN, keystore.StandardScryptP)
	if len(ks.Accounts()) == 0 {
		return fmt.Errorf("no accounts found in the key store")
	}
	account := ks.Accounts()[0]

	err := ks.Unlock(account, password)
	if err != nil {
		return fmt.Errorf("failed to unlock the account: %v", err)
	}

	fmt.Printf("Account %s loaded\n", account.Address.Hex())
	return nil
}
