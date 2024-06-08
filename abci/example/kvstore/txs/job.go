package txs

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// selectPseudorandomMiner selects a miner pseudorandomly based on the block height, app hash, and nonce.
func selectPseudorandomMiner(miners []string, blockHeight int, appHash string, nonce int) string {
	// Combine block height, app hash, and nonce into a single string
	combinedInput := fmt.Sprintf("%d%s%d", blockHeight, appHash, nonce)

	// Generate a SHA-256 hash of the combined input
	hash := sha256.New()
	hash.Write([]byte(combinedInput))
	hashBytes := hash.Sum(nil)

	// Convert the hash to a big integer
	hashInt := new(big.Int)
	hashInt.SetBytes(hashBytes)

	// Get the index within the range of the list of miners
	index := new(big.Int).Mod(hashInt, big.NewInt(int64(len(miners)))).Int64()

	// Select the miner at the generated index
	selectedMiner := miners[index]

	return selectedMiner
}

/*
func main() {
	// Example usage
	miners := []string{"miner1", "miner2", "miner3", "miner4"}
	blockHeight := 123456
	appHash := "exampleAppHash"
	nonce := 7890

	selectedMiner := selectPseudorandomMiner(miners, blockHeight, appHash, nonce)
	fmt.Println("Selected Miner:", selectedMiner)
}
*/
