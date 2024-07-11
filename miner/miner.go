package miner

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"os"
)

type Miner struct {
	RPCStatus   int // 0 for stale, 1 for active
	RPCEndpoint string
	Wallet      struct {
		Keystore *keystore.KeyStore
		Password string
	}
	KeyFilePath string
}

// NewMiner initializes a new Miner with the given RPC endpoint and keystore path
func NewMiner(rpcEndpoint, keyFilePath string) *Miner {
	m := &Miner{
		RPCStatus:   0,
		RPCEndpoint: rpcEndpoint,
		KeyFilePath: keyFilePath,
	}
	m.Wallet.Keystore = keystore.NewKeyStore(keyFilePath, keystore.StandardScryptN, keystore.StandardScryptP)
	return m
}

// Initialize checks or creates a key, sets up the wallet, and updates the RPC status
func (m *Miner) Initialize() error {
	if err := m.checkOrCreateKey(); err != nil {
		return err
	}
	return m.UpdateRPCStatus()
}

// UpdateRPCStatus checks the RPC status and updates the Miner struct
func (m *Miner) UpdateRPCStatus() error {
	err := QueryRPCStatus(m.RPCEndpoint)
	if err != nil {
		m.RPCStatus = 0
		return fmt.Errorf("RPC status check failed: %v", err)
	}
	m.RPCStatus = 1
	return nil
}

// checkOrCreateKey checks if the key file exists and either loads or creates a new key
func (m *Miner) checkOrCreateKey() error {
	_, err := os.Stat(m.KeyFilePath)
	if os.IsNotExist(err) {
		fmt.Println("No keyfile found. Creating a new key...")
		password, err := PromptPassword(true) // Ask for password with confirmation
		if err != nil {
			return fmt.Errorf("prompt for password failed: %v", err)
		}
		m.Wallet.Password = password
		_, err = CreateNewKey(m.KeyFilePath, password)
		if err != nil {
			return fmt.Errorf("failed to create new key: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check key file: %v", err)
	} else {
		fmt.Println("Keyfile found. Loading key...")
		password, err := PromptPassword(false) // Ask for password without confirmation
		if err != nil {
			return fmt.Errorf("prompt for password failed: %v", err)
		}
		m.Wallet.Password = password
		_, err = LoadKey(m.KeyFilePath, password)
		if err != nil {
			return fmt.Errorf("failed to load key: %v", err)
		}
	}
	return nil
}

// ToAddress returns the Ethereum address associated with the miner's primary account
func (m *Miner) ToAddress() (common.Address, error) {
	if len(m.Wallet.Keystore.Accounts()) == 0 {
		return common.Address{}, fmt.Errorf("no accounts in keystore")
	}
	account := m.Wallet.Keystore.Accounts()[0] // Get the first account
	return account.Address, nil
}

// ToAddressHex returns the Ethereum address in hexadecimal string format
func (m *Miner) ToAddressHex() (string, error) {
	address, err := m.ToAddress()
	if err != nil {
		return "", err
	}
	return address.Hex(), nil
}
