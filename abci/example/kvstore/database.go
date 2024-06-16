package kvstore

import (
	"encoding/json"
	"fmt"
	db "github.com/tendermint/tm-db"
)

type ClientInfo struct {
	Name  string
	Power uint64
}

// BuildKey generates a database key for a given Ethereum address.
func BuildKeyForClientRegistration(ethereumAddress string) []byte {
	return []byte(fmt.Sprintf("clientRegistration_%s", ethereumAddress))
}

// StoreClientInfo stores ClientInfo in the database under the key derived from the Ethereum address.
func StoreClientInfo(db db.DB, ethereumAddress string, info ClientInfo) error {
	dataBytes, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return db.Set(BuildKeyForClientRegistration(ethereumAddress), dataBytes)
}

// GetClientInfo retrieves ClientInfo from the database using the Ethereum address.
func GetClientInfo(db db.DB, ethereumAddress string) (ClientInfo, error) {
	dataBytes, err := db.Get(BuildKeyForClientRegistration(ethereumAddress))
	if err != nil {
		return ClientInfo{}, err
	}
	var info ClientInfo
	err = json.Unmarshal(dataBytes, &info)
	return info, err
}

type MinerInfo struct {
	Name         string   // The name of the miner
	Power        uint64   // The computational power of the miner, possibly in hashes per second
	ServiceTypes []uint64 // An array of service type identifiers that the miner provides
	IP           string   // The IP address of the miner for network connections
}

// BuildKeyForMinerRegistration generates a database key for a given Ethereum address.
func BuildKeyForMinerRegistration(ethereumAddress string) []byte {
	return []byte(fmt.Sprintf("minerRegistration_%s", ethereumAddress))
}

// StoreMinerInfo stores MinerInfo in the database under the key derived from the Ethereum address.
func StoreMinerInfo(db db.DB, ethereumAddress string, info MinerInfo) error {
	dataBytes, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return db.Set(BuildKeyForMinerRegistration(ethereumAddress), dataBytes)
}

// GetMinerInfo retrieves MinerInfo from the database using the Ethereum address.
func GetMinerInfo(db db.DB, ethereumAddress string) (MinerInfo, error) {
	dataBytes, err := db.Get(BuildKeyForMinerRegistration(ethereumAddress))
	if err != nil {
		return MinerInfo{}, err
	}
	var info MinerInfo
	err = json.Unmarshal(dataBytes, &info)
	return info, err
}

// BuildServiceTypeKey generates a database key for service type mappings.
func BuildServiceTypeKey(serviceType uint64) []byte {
	return []byte(fmt.Sprintf("serviceType_%d", serviceType))
}

// GetMinersForServiceType retrieves the list of miners for a given service type from the database.
func GetMinersForServiceType(db db.DB, serviceType uint64) ([]string, error) {
	key := BuildServiceTypeKey(serviceType)
	data, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return []string{}, nil // No miners registered yet for this service type.
	}

	var miners []string
	err = json.Unmarshal(data, &miners)
	if err != nil {
		return nil, err
	}
	return miners, nil
}

// StoreMinersForServiceType stores the list of miners for a given service type in the database.
func StoreMinersForServiceType(db db.DB, serviceType uint64, miners []string) error {
	key := BuildServiceTypeKey(serviceType)
	data, err := json.Marshal(miners)
	if err != nil {
		return err
	}
	return db.Set(key, data)
}

// RegisterMiner registers a new miner and updates the service type mappings in the database.
func RegisterMiner(db db.DB, miner MinerInfo, minerAddress string) error {
	for _, serviceType := range miner.ServiceTypes {
		miners, err := GetMinersForServiceType(db, serviceType)
		if err != nil {
			return fmt.Errorf("failed to get miners for service type %d: %v", serviceType, err)
		}

		// Check if the miner's address is already listed under this service type.
		found := false
		for _, addr := range miners {
			if addr == minerAddress {
				found = true
				break
			}
		}

		// If not found, add the miner's address to the list and store it back.
		if !found {
			miners = append(miners, minerAddress)
			err = StoreMinersForServiceType(db, serviceType, miners)
			if err != nil {
				return fmt.Errorf("failed to store miners for service type %d: %v", serviceType, err)
			}
		}
	}

	return nil
}
