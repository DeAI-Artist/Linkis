package txs

import "encoding/json"

type MessageContent interface {
	ToBytes() ([]byte, error)
}

const (
	ClientRegistrationType   = 1
	ServiceRequestType       = 2
	ClientRatingMsgType      = 3
	MinerRegistrationType    = 4
	MinerServiceDoneType     = 5
	MinerStatusUpdateType    = 6
	MinerRewardClaimType     = 7
	MinerServiceStartingType = 8
)

type ClientRegistrationMsg struct {
	ClientName string `json:"client_name"`
}

type ServiceRequestMsg struct {
	ServiceID uint64 `json:"service_id"`
	Meta      []byte `json:"meta"`
}

type ClientRatingMsg struct {
	ReviewedMinerAddr string `json:"miner_addr"`
	Rating            int    `json:"rating"`
}

type MinerRegistrationMsg struct {
	MinerName    string   `json:"miner_name"`    // The name of the miner
	ServiceTypes []uint64 `json:"service_types"` // An array of service type identifiers
	IP           string   `json:"ip"`            // The IP address of the miner for network connections
	Status       uint8    `json:"status"`
}

type MinerServiceDoneMsg struct {
	ServiceID   string `json:"service_id"`   // Unique identifier for the service
	ServiceType uint64 `json:"service_type"` // Numeric identifier of the type of service completed
}

type MinerStatusUpdateMsg struct {
	AddServiceTypes    []uint64 `json:"service_types"`
	RemoveServiceTypes []uint64 `json:"service_types"`
	Status             uint8    `json:"status"`
}

type MinerRewardClaimMsg struct{}

// ServiceStartingMsg defines the data structure for initiating a service process.
type ServiceStartingMsg struct {
	ServiceID       string `json:"service_id"`        // Unique identifier for the service
	MaxTimeoutBlock int64  `json:"max_timeout_block"` // Maximum number of blocks the service should run before automatic termination
}

// Implement ToBytes for each struct
func (m ClientRegistrationMsg) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

func (m ServiceRequestMsg) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

func (m ClientRatingMsg) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

func (m MinerRegistrationMsg) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

func (m MinerServiceDoneMsg) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

func (m MinerStatusUpdateMsg) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

func (m MinerRewardClaimMsg) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

// ToBytes converts the ServiceStartingMsg to a byte slice for easy transmission and storage.
func (m ServiceStartingMsg) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}
