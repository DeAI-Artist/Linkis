package kvstore

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestLoadWallet(t *testing.T) {
	// Generate a valid private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())

	// Test with a valid private key
	wallet, err := LoadWallet(privateKeyHex)
	if err != nil {
		t.Errorf("LoadWallet failed with valid private key: %v", err)
	}
	if wallet == nil {
		t.Errorf("LoadWallet returned nil for a valid private key")
	}

	// Test error handling with an invalid private key
	_, err = LoadWallet("invalidkey")
	if err == nil {
		t.Errorf("LoadWallet should fail with an invalid private key")
	}
}

func TestChainIDUsingSigner(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	chainID := big.NewInt(11155111) // Sepolia Chain ID

	// Mock transaction to test signing
	tx := types.NewTransaction(0, common.HexToAddress("0xAddress"), big.NewInt(0), 21000, big.NewInt(1), nil)
	signer := types.LatestSignerForChainID(chainID)

	// Sign the transaction
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	// Check if the recovered ChainID matches
	if signedTx.ChainId().Cmp(chainID) != 0 {
		t.Errorf("Expected Chain ID %v, got %v", chainID, signedTx.ChainId())
	}
}
