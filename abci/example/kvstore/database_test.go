package kvstore

import (
	"encoding/json"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

// TestStoreAndGetClientInfo tests the storing and retrieval of ClientInfo in the database.
func TestStoreAndGetClientInfo(t *testing.T) {
	// Initialize the in-memory database
	db := dbm.NewMemDB()

	// Example client information
	clientInfo := ClientInfo{
		Name:  "Alice",
		Power: 10,
	}
	ethereumAddress := "0x123abc"

	// Store client information
	err := StoreClientInfo(db, ethereumAddress, clientInfo)
	if err != nil {
		t.Fatalf("StoreClientInfo failed: %s", err)
	}

	// Retrieve client information
	retrievedInfo, err := GetClientInfo(db, ethereumAddress)
	if err != nil {
		t.Fatalf("GetClientInfo failed: %s", err)
	}

	// Compare the stored and retrieved information
	if retrievedInfo.Name != clientInfo.Name || retrievedInfo.Power != clientInfo.Power {
		t.Errorf("Retrieved info does not match stored info. Got %+v, want %+v", retrievedInfo, clientInfo)
	}
}

func TestMapInMemDB(t *testing.T) {
	// Initialize an in-memory database
	memDB := dbm.NewMemDB()

	// Original map
	originalMap := map[string]int{"a": 1, "b": 2, "c": 3}

	// Serialize the map to JSON for storage
	jsonData, err := json.Marshal(originalMap)
	if err != nil {
		t.Fatalf("Failed to serialize map: %v", err)
	}

	// Store serialized data
	if err := memDB.Set([]byte("mapKey"), jsonData); err != nil {
		t.Fatalf("Failed to set data in memDB: %v", err)
	}

	// Retrieve and deserialize the map
	retrievedData, err := memDB.Get([]byte("mapKey"))
	if err != nil {
		t.Fatalf("Failed to get data from memDB: %v", err)
	}
	retrievedMap := make(map[string]int)
	if err := json.Unmarshal(retrievedData, &retrievedMap); err != nil {
		t.Fatalf("Failed to deserialize map: %v", err)
	}

	// Compare the original and retrieved maps
	if len(retrievedMap) != len(originalMap) {
		t.Fatalf("Mismatch in map sizes: expected %v, got %v", len(originalMap), len(retrievedMap))
	}
	for key, expectedValue := range originalMap {
		if value, exists := retrievedMap[key]; !exists || value != expectedValue {
			t.Fatalf("Mismatch for key %s: expected %v, got %v", key, expectedValue, value)
		}
	}
}

func TestStoreAndGetMinerInfo(t *testing.T) {
	// Initialize the in-memory database
	db := dbm.NewMemDB()

	// Example miner information
	minerInfo := MinerInfo{
		Name:         "Miner Bob",
		Power:        500,
		ServiceTypes: []uint64{1, 2, 3},
		IP:           "192.168.1.100:8080", // IP address and port combined in one string
	}
	ethereumAddress := "0x456def"

	// Store miner information
	err := StoreMinerInfo(db, ethereumAddress, minerInfo)
	if err != nil {
		t.Fatalf("StoreMinerInfo failed: %s", err)
	}

	// Retrieve miner information
	retrievedInfo, err := GetMinerInfo(db, ethereumAddress)
	if err != nil {
		t.Fatalf("GetMinerInfo failed: %s", err)
	}

	// Compare the stored and retrieved information
	if retrievedInfo.Name != minerInfo.Name || retrievedInfo.Power != minerInfo.Power ||
		len(retrievedInfo.ServiceTypes) != len(minerInfo.ServiceTypes) || retrievedInfo.IP != minerInfo.IP {
		t.Errorf("Retrieved info does not match stored info. Got %+v, want %+v", retrievedInfo, minerInfo)
	}

	// Optionally check if all service types match
	for i, serviceType := range retrievedInfo.ServiceTypes {
		if serviceType != minerInfo.ServiceTypes[i] {
			t.Errorf("Service type mismatch at index %d: got %d, want %d", i, serviceType, minerInfo.ServiceTypes[i])
		}
	}
}
