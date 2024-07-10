package commands

import (
	"fmt"
	"github.com/DeAI-Artist/MintAI/miner" // Adjust this import path to where your wallet.go resides
	"github.com/spf13/cobra"
	"os"
)

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
		password := miner.PromptPassword(true) // Ask for password with confirmation
		if err := miner.CreateNewKey(keyfilePath, password); err != nil {
			return fmt.Errorf("failed to create new key: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("Failed to check key file: %v", err)
	} else {
		fmt.Println("Keyfile found. Loading key...")
		password := miner.PromptPassword(false) // Ask for password without confirmation
		if err := miner.LoadKey(keyfilePath, password); err != nil {
			return fmt.Errorf("failed to load key: %v", err)
		}
	}

	return nil
}
