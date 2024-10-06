package miner

import (
	"bufio"
	"encoding/json"
	"fmt"
	kv "github.com/DeAI-Artist/Linkis/abci/example/kvstore"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
	"os"
	"strings"
)

// MinerRegistrationRequest represents the data required to register a miner.
type MinerRegistrationRequest struct {
	MinerName    string   `json:"miner_name"`    // The name of the miner
	ServiceTypes []uint64 `json:"service_types"` // An array of service type identifiers
	ServicePort  uint64   `json:"service_port"`  // The network port the miner uses
	IP           string   `json:"ip"`            // The IP address of the miner
	Status       uint8    `json:"status"`        // The operational status of the miner
}

type Miner struct {
	RPCStatus   int // 0 for stale, 1 for active
	RPCEndpoint string
	Wallet      struct {
		Keystore *keystore.KeyStore
		Password string
	}
	KeyFilePath  string
	IP           string   `json:"ip"`
	Status       uint8    `json:"status"`
	MinerName    string   `json:"miner_name"`
	ServiceTypes []uint64 `json:"service_types"`
	ServicePort  uint64   `json:"service_port"`
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
	// Check or create a key for the miner
	if err := m.checkOrCreateKey(); err != nil {
		return fmt.Errorf("failed to check or create key: %v", err)
	}

	// Update the RPC status to ensure it is current before registration
	if err := m.UpdateRPCStatus(); err != nil {
		return fmt.Errorf("failed to update RPC status: %v", err)
	}

	// Register the miner if not already registered
	if err := m.RegisterMiner(); err != nil {
		return fmt.Errorf("failed to register miner: %v", err)
	}

	return nil
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

// ToAddress returns the Ethereum address associated with the miner's primary account.
// It creates a new account if none exist in the keystore.
func (m *Miner) ToAddress() (common.Address, error) {
	// Check if there are any accounts in the keystore
	if len(m.Wallet.Keystore.Accounts()) == 0 {
		fmt.Println("No accounts found in keystore, creating a new key...")
		account, err := m.createNewKey()
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to create new key: %v", err)
		}
		fmt.Println("New key created with address:", account.Address.Hex())
		return account.Address, nil
	}

	// If accounts are available, return the address of the first account
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

func (m *Miner) RegisterMiner() error {
	// First check if the miner is already registered
	address, err := m.ToAddressHex()
	if err != nil {
		return fmt.Errorf("error getting miner's Ethereum address: %v", err)
	}
	registered, err := IsMinerRegistered(m.RPCEndpoint, address)
	if err != nil {
		return fmt.Errorf("error checking if miner is registered: %v", err)
	}
	if registered {
		fmt.Println("Miner is already registered.")
		// TODO: Handle already registered miner
		return nil
	}

	scanner := bufio.NewScanner(os.Stdin)

	// Prompt for miner name if not set
	if m.MinerName == "" {
		fmt.Print("Enter miner name: ")
		scanner.Scan()
		m.MinerName = scanner.Text()
	}

	// Prompt for service types with option for default
	if len(m.ServiceTypes) == 0 {
		prompt := promptui.Select{
			Label: "Do you want to use the default service types?",
			Items: []string{"Yes", "No"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			return fmt.Errorf("prompt failed %v", err)
		}

		if strings.ToLower(result) == "yes" {
			defaultTypes, err := GetSystemServiceTypes(m.RPCEndpoint)
			if err != nil {
				return fmt.Errorf("failed to get default service types: %v", err)
			}
			m.ServiceTypes = defaultTypes
		} else {
			fmt.Print("Enter service types (comma-separated): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			serviceTypesInput := scanner.Text()
			for _, st := range strings.Split(serviceTypesInput, ",") {
				var serviceType uint64
				fmt.Sscanf(st, "%d", &serviceType)
				m.ServiceTypes = append(m.ServiceTypes, serviceType)
			}
		}
	}

	// Prompt for service port if not set
	if m.ServicePort == 0 {
		fmt.Print("Enter service port (press Enter for default): ")
		scanner.Scan()
		servicePortInput := scanner.Text()
		if servicePortInput == "" {
			m.ServicePort = 26688 // default port
		} else {
			fmt.Sscanf(servicePortInput, "%d", &m.ServicePort)
		}
	}

	// Assuming IP is set dynamically
	if m.IP == "" {
		var err error
		m.IP, err = GetPublicIP()
		if err != nil {
			// Handle the error appropriately
			fmt.Errorf("failed to fetch public IP: %v", err)
		} else {
			fmt.Println("Public IP set to:", m.IP)
		}
	}
	if m.Status == kv.Stale {
		m.Status = kv.Ready // Example static status
	}

	// Construct the request payload
	requestPayload := MinerRegistrationRequest{
		MinerName:    m.MinerName,
		ServiceTypes: m.ServiceTypes,
		ServicePort:  m.ServicePort,
		IP:           m.IP,
		Status:       m.Status,
	}

	// Marshal the request into JSON
	requestData, err := json.Marshal(requestPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %v", err)
	}
	spew.Dump("request payload: ", requestData)

	return nil
}

// createNewKey is a helper method to create a new key and store it in the keystore
func (m *Miner) createNewKey() (accounts.Account, error) {
	// Use the password from the Miner struct to create a new key
	if m.Wallet.Password == "" {
		return accounts.Account{}, fmt.Errorf("password not set in Miner struct")
	}
	ks := keystore.NewKeyStore(m.KeyFilePath, keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(m.Wallet.Password)
	if err != nil {
		return accounts.Account{}, fmt.Errorf("failed to create new key: %v", err)
	}

	// Optionally, save the keystore changes to disk or handle them appropriately
	// err = ks.StoreKey(m.KeyFilePath, account, m.Wallet.Password)
	// if err != nil {
	//     return accounts.Account{}, fmt.Errorf("failed to save the key to disk: %v", err)
	// }

	return account, nil
}
