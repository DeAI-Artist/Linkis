package kvstore

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// LoadWallet initializes a wallet with an EVM-compatible private key for transaction signing on the Sepolia network.
func LoadWallet(privateKeyHex string) (*bind.TransactOpts, error) {
	if privateKeyHex == "" {
		return nil, errors.New("private key cannot be empty")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	chainID := big.NewInt(11155111) // Sepolia chain ID
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
