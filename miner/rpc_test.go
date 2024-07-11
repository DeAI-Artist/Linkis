package miner

import (
	"github.com/jarcoal/httpmock"
	"testing"
)

// TestQueryRPCStatus checks both success and failure cases for the QueryRPCStatus function
func TestQueryRPCStatus(t *testing.T) {
	// Activate the HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Setup the responses for both test cases
	successResponder := httpmock.NewStringResponder(200, `{"jsonrpc":"2.0","id":-1,"result":{"node_info":{"protocol_version":{"p2p":"8","block":"11","app":"1"},"id":"d8629baf3b637dd87ab890c2ebee5f589f17af05","listen_addr":"tcp://0.0.0.0:26656","network":"chain-rYTYBC","version":"main-80a7e77f053b148bd890d5f635a0aba8b1dcf86b","channels":"40202122233038606100","moniker":"134.209.85.94","other":{"tx_index":"on","rpc_address":"tcp://0.0.0.0:26657"}},"sync_info":{"latest_block_hash":"507D282D72B5AA9B33DDF28205A07BBDE1450D1503856F132617B6D672D32597","latest_app_hash":"AFA9F60BB41DFB624A51F31437B44E284758DEF3DE892C2CD424A37C5912708F","latest_block_height":"196243","latest_block_time":"2024-07-11T09:38:00.327118298Z","earliest_block_hash":"F9B2DE42FE03D0EF570651089C83A52BBA613173A14E1A8F25BCEB92901B0818","earliest_app_hash":"","earliest_block_height":"1","earliest_block_time":"2024-07-07T17:28:24.588464133Z","catching_up":false},"validator_info":{"address":"FFDD18B0A49FDCC65EF9AD1B72102251113F04D1","pub_key":{"type":"tendermint/PubKeyEd25519","value":"4OpftiHIIbk9Ii/znxd889LhJvcCXi6zC6twVAPp6Ho="},"voting_power":"1"}}}`)
	failureResponder := httpmock.NewStringResponder(500, `Internal Server Error`)

	httpmock.RegisterResponder("GET", "http://134.209.85.94:26657/status", successResponder)
	httpmock.RegisterResponder("GET", "http://localhost:26657/status", failureResponder)

	// Test success case
	err := QueryRPCStatus("134.209.85.94:26657")
	if err != nil {
		t.Errorf("Expected no error for 134.209.85.94:26657, got %v", err)
	}

	// Test failure case
	err = QueryRPCStatus("localhost:26657")
	if err == nil {
		t.Errorf("Expected an error for localhost:26657, got nil")
	}
}
