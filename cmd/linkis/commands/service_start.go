package commands

import (
	"fmt"
	"github.com/DeAI-Artist/Linkis/miner"
	"github.com/spf13/cobra"
	"log"
)

var rpcEndpoint string

var ServiceStart = &cobra.Command{
	Use:     "service-start",
	Aliases: []string{"service_start"},
	Short:   "Generate task commitment to the network",
	PreRun:  deprecateSnakeCase,
	RunE:    serviceStartWithConfig,
}

func init() {
	// Define the flag for the RPC endpoint
	ServiceStart.Flags().StringVarP(&rpcEndpoint, "rpc-endpoint", "r", "http://localhost:8545", "RPC endpoint URL")
}

func serviceStartWithConfig(cmd *cobra.Command, args []string) error {
	fmt.Printf("Using RPC Endpoint: %s\n", rpcEndpoint) // Example usage of the RPC endpoint

	// Here you fetch the key file path from your configuration, adjust as necessary
	keyfilePath := config.NodeMinerKeyFile() // Assuming there's a function in config package

	// Create a new Miner instance with the specified RPC endpoint and key file path
	minerInstance := miner.NewMiner(rpcEndpoint, keyfilePath)

	// Initialize the miner (checks or creates a key, sets up the wallet, updates RPC status)
	if err := minerInstance.Initialize(); err != nil {
		log.Fatalf("Failed to initialize miner: %v", err)
		return err
	}

	fmt.Println("Miner initialized and ready.")
	return nil
}

/*
# Basic syntax
linkis service-start --rpc-endpoint "http://your-rpc-server:port"

# Example usage
linkis service-start --rpc-endpoint "http://localhost:8545"

# Or using the shorthand for the flag
linkis service-start -r "http://localhost:8545"
*/
