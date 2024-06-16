package txs

import "encoding/json"

type MessageContent interface {
	ToBytes() ([]byte, error)
}

const (
	ClientRegistrationType = 1
	ServiceRequestType     = 2
	ClientRatingMsgType    = 3
	MinerRegistrationType  = 4
	MinerServiceDoneType   = 5
	MinerStatusUpdateType  = 6
	MinerRewardClaimType   = 7
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
}

type MinerServiceDoneMsg struct {
	ServiceID string `json:"service_id"`
}

type MinerStatusUpdateMsg struct {
	AddServiceTypes    []uint64 `json:"service_types"`
	RemoveServiceTypes []uint64 `json:"service_types"`
	Status             string   `json:"status"`
}

type MinerRewardClaimMsg struct{}

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
