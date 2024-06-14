package kvstore

import (
	"testing"

	dbm "github.com/tendermint/tm-db"
)

// TestStoreAndGetClientInfo tests the storing and retrieval of ClientInfo in the database.
func TestStoreAndGetClientInfo(t *testing.T) {
	// Initialize the in-memory database
	db := dbm.NewMemDB()

	// Example client information
	clientInfo := ClientInfo{
		Name:  "Alice",
		power: 30,
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
	if retrievedInfo.Name != clientInfo.Name || retrievedInfo.power != clientInfo.power {
		t.Errorf("Retrieved info does not match stored info. Got %+v, want %+v", retrievedInfo, clientInfo)
	}
}
