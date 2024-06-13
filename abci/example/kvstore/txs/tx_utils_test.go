package txs

import (
	"testing"
)

func TestHashAndRecover(t *testing.T) {
	message := []byte("Example `personal_sign` message")
	expectedAddress := "0x0921e25C5CDd1853c2CF3014D4B0b46d0B2C31ab"
	sigS := "0x7635254a22cd7762ded328cdb27292884ee2ea21500bbff52ff549b09ef0aaa80dc5e5bc03d67c3ce5176747dcec576bf64d3959bf4a5b0b890fa188026825531b"
	signature, err := HexToBytes(sigS[2:]) // remove the "0x" prefix
	if err != nil {
		t.Fatalf("Failed to decode signature: %s", err)
	}
	hashed := HashPersonalMessage(message)
	// Sample signature (this needs to be a real one for actual tests)

	pubKey, err := RecoverPubKey(hashed, signature)
	if err != nil {
		t.Fatalf("Failed to recover public key: %s", err)
	}

	address := AddressFromPublicKey(pubKey)
	if address != expectedAddress {
		t.Errorf("Expected address %s, got %s", expectedAddress, address)
	}
}
