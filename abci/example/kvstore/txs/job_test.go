package txs

import (
	"testing"
)

// TestSelectPseudorandomMiner checks the functionality of selectPseudorandomMiner by verifying consistency and variability of output.
func TestSelectPseudorandomMiner(t *testing.T) {
	miners := []string{
		"Miner1", "Miner2", "Miner3", "Miner4", "Miner5",
		"Miner6", "Miner7", "Miner8", "Miner9", "Miner10",
		"Miner11", "Miner12", "Miner13", "Miner14", "Miner15",
		"Miner16", "Miner17", "Miner18", "Miner19", "Miner20",
		"Miner21", "Miner22", "Miner23", "Miner24", "Miner25",
		"Miner26", "Miner27", "Miner28", "Miner29", "Miner30",
	}

	// Test for consistency: same inputs should always yield the same miner
	result1 := SelectPseudorandomMiner(miners, 100, "abc123", "service1")
	result2 := SelectPseudorandomMiner(miners, 100, "abc123", "service1")
	if result1 != result2 {
		t.Errorf("Expected consistent result but got %s and %s", result1, result2)
	}

	// Test for variability: different inputs should yield different results
	result3 := SelectPseudorandomMiner(miners, 101, "abc123", "service1") // Changing blockHeight
	result4 := SelectPseudorandomMiner(miners, 100, "abc124", "service1") // Changing appHash
	result5 := SelectPseudorandomMiner(miners, 100, "abc123", "service2") // Changing serviceID
	//spew.Dump(result3, result4, result5)
	// Check if any of the results are the same when they shouldn't be
	if result1 == result3 || result1 == result4 || result1 == result5 {
		t.Errorf("Expected different results but got same miner. Results: %s, %s, %s, %s", result1, result3, result4, result5)
	}
}
