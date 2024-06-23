package kvstore

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	dbm "github.com/tendermint/tm-db"
	"reflect"
	"strconv"
	"testing"
	"time"
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

func TestRegisterMiner(t *testing.T) {
	db := dbm.NewMemDB()

	minerInfo := MinerInfo{
		Name:         "Miner Bob",
		Power:        500,
		ServiceTypes: []uint64{1, 2},
		IP:           "192.168.1.100:8080",
	}
	minerAddress := "0x456def"

	// Case 1: No service type key exists yet
	err := RegisterMiner(db, minerInfo, minerAddress)
	if err != nil {
		t.Fatalf("RegisterMiner failed: %s", err)
	}
	for _, serviceType := range minerInfo.ServiceTypes {
		miners, _ := GetMinersForServiceType(db, serviceType)
		if len(miners) != 1 || miners[0] != minerAddress {
			t.Errorf("Expected [%s] for service type %d, got %v", minerAddress, serviceType, miners)
		}
	}

	// Case 2: Service type key exists and already includes the specified miner
	err = RegisterMiner(db, minerInfo, minerAddress)
	if err != nil {
		t.Fatalf("RegisterMiner failed on second call: %s", err)
	}
	for _, serviceType := range minerInfo.ServiceTypes {
		miners, _ := GetMinersForServiceType(db, serviceType)
		if len(miners) != 1 || miners[0] != minerAddress { // No duplicate should be added
			t.Errorf("Unexpected duplicate or change in miners list for service type %d, got %v", serviceType, miners)
		}
	}

	// Case 3: Service type key exists but does not include the specified miner
	newMinerAddress := "0x123abc"
	err = RegisterMiner(db, minerInfo, newMinerAddress)
	if err != nil {
		t.Fatalf("RegisterMiner failed on third call: %s", err)
	}
	for _, serviceType := range minerInfo.ServiceTypes {
		miners, _ := GetMinersForServiceType(db, serviceType)
		if !reflect.DeepEqual(miners, []string{minerAddress, newMinerAddress}) {
			t.Errorf("Expected [%s, %s] for service type %d, got %v", minerAddress, newMinerAddress, serviceType, miners)
		}
	}

	// Additional case: Database operations error handling
	// You could mock the database to throw errors and verify that RegisterMiner handles it properly.
	// This would typically require an interface for the database and a mock implementation.
}

// TestMinerStatusManagement tests adding, updating, querying, and removing miner statuses.
func TestMinerStatusManagement(t *testing.T) {
	db := dbm.NewMemDB()
	address := "0x123"
	initialStatus := Ready
	updatedStatus := Busy

	// Add a new miner
	if err := AddOrUpdateMinerStatus(db, address, initialStatus); err != nil {
		t.Fatalf("Failed to add new miner status: %v", err)
	}

	// Verify initial status
	status, err := GetMinerStatus(db, address)
	if err != nil || status != initialStatus {
		t.Errorf("Expected initial status %d for address %s; got %d", initialStatus, address, status)
	}

	// Update the miner's status
	if err := AddOrUpdateMinerStatus(db, address, updatedStatus); err != nil {
		t.Errorf("Failed to update miner status: %v", err)
	}

	// Verify updated status
	status, err = GetMinerStatus(db, address)
	if err != nil || status != updatedStatus {
		t.Errorf("Expected updated status %d for address %s; got %d", updatedStatus, address, status)
	}

	// Remove the miner's status
	if err := RemoveMinerStatus(db, address); err != nil {
		t.Errorf("Failed to remove miner status: %v", err)
	}

	// Verify removal
	_, err = GetMinerStatus(db, address)
	if err == nil {
		t.Errorf("Expected error for fetching status of removed miner, got none")
	}
}

// BenchmarkStatusOperations benchmarks adding, retrieving, and removing miner statuses.
func BenchmarkStatusOperations(b *testing.B) {
	db := dbm.NewMemDB()
	for n := 1; n <= b.N; n++ {
		address := "miner" + strconv.Itoa(n) // Create unique miner address
		status := uint8(n % 3)               // Cycle through 0, 1, 2 statuses

		// Add or update miner status
		b.Run("AddOrUpdate", func(b *testing.B) {
			if err := AddOrUpdateMinerStatus(db, address, status); err != nil {
				b.Error("Failed to add/update miner status:", err)
			}
		})

		// Get miner status
		b.Run("GetStatus", func(b *testing.B) {
			if _, err := GetMinerStatus(db, address); err != nil {
				b.Error("Failed to get miner status:", err)
			}
		})

		// Remove miner status
		b.Run("RemoveStatus", func(b *testing.B) {
			if err := RemoveMinerStatus(db, address); err != nil {
				b.Error("Failed to remove miner status:", err)
			}
		})
	}
}

func PrepopulateDatabase(db dbm.DB) error {
	statuses := make(MinerStatuses) // Create a new map to store the statuses

	// Populate the statuses map in memory
	for i := 0; i < 100000; i++ {
		address := "miner" + strconv.Itoa(i)
		status := uint8(i % 3)     // Cycle statuses between 0, 1, and 2
		statuses[address] = status // Directly update the map
	}

	// Save the updated map back to the database only once after all updates
	if err := SaveMinerStatuses(db, statuses); err != nil {
		return fmt.Errorf("failed to save miner statuses: %v", err)
	}
	return nil
}

// TestDatabaseOperations tests add, remove, update, and query operations on a large dataset.
func TestDatabaseOperations(t *testing.T) {
	db := dbm.NewMemDB()
	// Populate the database with initial data
	if err := PrepopulateDatabase(db); err != nil {
		t.Fatalf("Error prepopulating database: %v", err)
	}
	println("data generation complete")

	// Start testing operations
	t.Run("AddOperation", func(t *testing.T) {
		// Measure time taken to add a new status
		start := time.Now()
		err := AddOrUpdateMinerStatus(db, "miner10001", 1)
		duration := time.Since(start)
		if err != nil {
			t.Errorf("Failed to add new miner status: %v", err)
		}
		t.Logf("Time taken to add new miner status: %v", duration)
	})

	t.Run("UpdateOperation", func(t *testing.T) {
		// Measure time taken to update an existing status
		start := time.Now()
		err := AddOrUpdateMinerStatus(db, "miner500000", 2) // Update miner at halfway point
		duration := time.Since(start)
		if err != nil {
			t.Errorf("Failed to update miner status: %v", err)
		}
		t.Logf("Time taken to update miner status: %v", duration)
	})

	t.Run("RemoveOperation", func(t *testing.T) {
		// Measure time taken to remove a status
		start := time.Now()
		err := RemoveMinerStatus(db, "miner500000")
		duration := time.Since(start)
		if err != nil {
			t.Errorf("Failed to remove miner status: %v", err)
		}
		t.Logf("Time taken to remove miner status: %v", duration)
	})

	t.Run("QueryOperation", func(t *testing.T) {
		// Measure time taken to query a status
		start := time.Now()
		_, err := GetMinerStatus(db, "miner50000")
		duration := time.Since(start)
		if err != nil {
			t.Errorf("Failed to query miner status: %v", err)
		}
		t.Logf("Time taken to query miner status: %v", duration)
	})
}

// TestAddMinerToServiceTypeMapping tests the addition of a miner to a service type.
func TestAddMinerToServiceTypeMapping(t *testing.T) {
	db := dbm.NewMemDB() // Using Tendermint's in-memory DB for testing
	serviceType := uint64(1)
	minerAddr := "miner123"

	// Initially no miners registered for this service type
	miners, _ := GetMinersForServiceType(db, serviceType)
	assert.Equal(t, 0, len(miners), "Service type should initially have no miners")

	// Add miner and check
	err := AddMinerToServiceTypeMapping(db, serviceType, minerAddr)
	assert.NoError(t, err, "Adding miner to service type should not produce an error")

	miners, _ = GetMinersForServiceType(db, serviceType)
	assert.Equal(t, 1, len(miners), "Service type should have one miner after adding")
	assert.Equal(t, minerAddr, miners[0], "The added miner address should match")
}

// TestRemoveMinerFromServiceTypeMapping tests the removal of a miner from a service type.
func TestRemoveMinerFromServiceTypeMapping(t *testing.T) {
	db := dbm.NewMemDB() // Using Tendermint's in-memory DB for testing
	serviceType := uint64(1)
	minerAddr := "miner123"

	// Add a miner first
	_ = AddMinerToServiceTypeMapping(db, serviceType, minerAddr)

	// Check miner is added
	miners, _ := GetMinersForServiceType(db, serviceType)
	assert.Equal(t, 1, len(miners), "Service type should have one miner before removal")

	// Remove miner and check
	err := RemoveMinerFromServiceTypeMapping(db, serviceType, minerAddr)
	assert.NoError(t, err, "Removing miner from service type should not produce an error")

	miners, _ = GetMinersForServiceType(db, serviceType)
	assert.Equal(t, 0, len(miners), "Service type should have no miners after removal")
}

// TestStoreAndGetClientRating tests the storing and retrieving functionality of client ratings.
func TestStoreAndGetClientRating(t *testing.T) {
	db := dbm.NewMemDB() // Initialize in-memory database
	minerAddress := "miner1"

	// Test case 1: Initial store and retrieve
	initialRatings := map[string]uint8{"client1": 5, "client2": 4}
	err := StoreClientRating(db, minerAddress, initialRatings)
	assert.NoError(t, err, "Storing initial ratings should not produce an error")

	retrievedRatings, err := GetClientRating(db, minerAddress)
	assert.NoError(t, err, "Retrieving initial ratings should not produce an error")
	assert.Equal(t, initialRatings, retrievedRatings, "Retrieved initial ratings should match the stored ratings")

	// Test case 2: Update an existing rating
	updatedRatings := map[string]uint8{"client1": 3}
	err = StoreClientRating(db, minerAddress, updatedRatings)
	assert.NoError(t, err, "Updating ratings should not produce an error")

	retrievedRatings, err = GetClientRating(db, minerAddress)
	assert.NoError(t, err, "Retrieving updated ratings should not produce an error")
	assert.Equal(t, updatedRatings, retrievedRatings, "Retrieved ratings after update should match the stored ratings")

	// Test case 3: Add a new rating
	updatedRatings["client3"] = 5
	err = StoreClientRating(db, minerAddress, updatedRatings)
	assert.NoError(t, err, "Adding a new rating should not produce an error")

	retrievedRatings, err = GetClientRating(db, minerAddress)
	assert.NoError(t, err, "Retrieving ratings after adding a new client should not produce an error")
	assert.Equal(t, updatedRatings, retrievedRatings, "Retrieved ratings after adding a new client should match the stored ratings")

	// Test case 4: Empty ratings
	emptyRatings := make(map[string]uint8)
	err = StoreClientRating(db, "miner2", emptyRatings)
	assert.NoError(t, err, "Storing empty ratings should not produce an error")

	retrievedRatings, err = GetClientRating(db, "miner2")
	assert.NoError(t, err, "Retrieving empty ratings should not produce an error")
	assert.Equal(t, emptyRatings, retrievedRatings, "Retrieved empty ratings should be empty")
}

func TestGenerateHashForMinerInfo(t *testing.T) {
	minerAddress := "miner123"
	metadata := []byte("test data")
	blockHeight := int64(12345)

	// Generate hash
	hash := GenerateHashForServiceInfo(minerAddress, metadata, blockHeight)
	//spew.Dump(hash)
	// Ensure hash is not empty
	assert.NotEmpty(t, hash, "The generated hash should not be empty")

	// Test for changes in any input resulting in a different hash
	differentMetadata := []byte("different data")

	differentHash := GenerateHashForServiceInfo(minerAddress, differentMetadata, blockHeight)
	//spew.Dump(differentHash)
	assert.NotEqual(t, hash, differentHash, "Hashes should differ with different metadata")
}

func TestJobInfoStorageAndRetrieval(t *testing.T) {
	db := dbm.NewMemDB()
	minerID := "miner123"
	job1 := JobInfo{
		ServiceID:   "service123",
		ClientID:    "client456",
		ServiceType: 789,
		JobStatus:   Registered, // Assuming Ready is a predefined constant
	}

	job2 := JobInfo{
		ServiceID:   "service123", // Same ServiceID to test update functionality
		ClientID:    "client789",
		ServiceType: 790,
		JobStatus:   Processing, // Different status to verify update mechanism
	}

	// Store the first job info
	assert.NoError(t, StoreJobInfo(db, minerID, job1), "Storing job info should not produce an error")

	// Retrieve the jobs list
	retrievedJobs, err := GetJobInfos(db, minerID)
	assert.NoError(t, err, "Retrieving job infos should not produce an error")
	assert.Len(t, retrievedJobs, 1, "There should be one job info stored initially")
	assert.Equal(t, job1, retrievedJobs[0], "The retrieved job info should match the initial stored info")

	// Store the second job info which has the same ServiceID to test the update functionality
	assert.NoError(t, StoreJobInfo(db, minerID, job2), "Updating job info should not produce an error")

	// Retrieve the jobs list again to check updates
	updatedJobs, err := GetJobInfos(db, minerID)
	assert.NoError(t, err, "Retrieving job infos after update should not produce an error")
	assert.Len(t, updatedJobs, 1, "There should still be only one job info after update")
	assert.Equal(t, job2, updatedJobs[0], "The retrieved job info should match the updated info")
}

func TestServiceRequestUtilities(t *testing.T) {
	// Initialize the in-memory database
	memDB := dbm.NewMemDB()

	// Test adding a new service request
	err := AddServiceRequest(memDB, "service1234", "miner5678", 102)
	assert.Nil(t, err, "AddServiceRequest should not return an error")

	// Verify that the service request has been added correctly
	requests, err := LoadServiceRequests(memDB)
	assert.Nil(t, err, "LoadServiceRequests should not return an error")
	assert.Equal(t, 1, len(requests), "There should be one service request")
	assert.Equal(t, "service1234", requests[0].ServiceID, "ServiceID should match")
	assert.Equal(t, "miner5678", requests[0].MinerID, "MinerID should match")
	assert.Equal(t, int64(102), requests[0].Height, "Height should match")

	// Remove the added service request by retaining above its height
	err = RetainServiceRequestsAboveHeight(memDB, 102)
	assert.Nil(t, err, "RetainServiceRequestsAboveHeight should not return an error")

	// Verify removal
	requests, err = LoadServiceRequests(memDB)
	assert.Nil(t, err, "LoadServiceRequests should not return an error after removal")
	assert.Equal(t, 0, len(requests), "All requests should be removed")

	// Add multiple new service requests
	heights := []int64{100, 101, 102, 103, 103, 104, 105, 105, 105, 106}
	for i, height := range heights {
		err = AddServiceRequest(memDB, fmt.Sprintf("service%d", i), fmt.Sprintf("miner%d", i), height)
		assert.Nil(t, err, "AddServiceRequest should not return an error when adding multiple")
	}

	// Verify all are added
	requests, err = LoadServiceRequests(memDB)
	assert.Nil(t, err, "LoadServiceRequests should not return an error after multiple adds")
	assert.Equal(t, 10, len(requests), "There should be ten service requests added")

	// Retain requests above a height of 105, which is repeated
	err = RetainServiceRequestsAboveHeight(memDB, 105)
	assert.Nil(t, err, "RetainServiceRequestsAboveHeight should not return an error on height 105")

	// Verify retention
	requests, err = LoadServiceRequests(memDB)
	assert.Nil(t, err, "LoadServiceRequests should not return an error after retention")
	assert.Equal(t, 1, len(requests), "Only one request should remain")
	assert.Equal(t, "service9", requests[0].ServiceID, "The remaining request should have the highest height")

	// Test retention below the smallest height in the list
	err = RetainServiceRequestsAboveHeight(memDB, 99)
	assert.Nil(t, err, "RetainServiceRequestsAboveHeight should not return an error on height below min")
	requests, err = LoadServiceRequests(memDB)
	assert.Nil(t, err, "LoadServiceRequests should not return an error after extreme low retention")
	assert.Equal(t, 1, len(requests), "All requests should remain when retaining below the smallest height")

	// Re-add and test retention above the maximum height in the list
	for i, height := range heights {
		err = AddServiceRequest(memDB, fmt.Sprintf("service%d", i), fmt.Sprintf("miner%d", i), height)
		assert.Nil(t, err, "Re-adding service requests should not return an error")
	}

	err = RetainServiceRequestsAboveHeight(memDB, 107)
	assert.Nil(t, err, "RetainServiceRequestsAboveHeight should not return an error on height above max")
	requests, err = LoadServiceRequests(memDB)
	assert.Nil(t, err, "LoadServiceRequests should not return an error after extreme high retention")
	assert.Equal(t, 0, len(requests), "No requests should remain when retaining above the highest height")
}
