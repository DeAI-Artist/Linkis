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

// TestHexToBytes tests the HexToBytes function with and without "0x" prefix.
func TestHexToBytes(t *testing.T) {
	tests := []struct {
		name     string
		hexStr   string
		expected []byte
		wantErr  bool
	}{
		{
			name:     "with 0x prefix",
			hexStr:   "0x1a2b3c",
			expected: []byte{0x1a, 0x2b, 0x3c},
			wantErr:  false,
		},
		{
			name:     "without 0x prefix",
			hexStr:   "1a2b3c",
			expected: []byte{0x1a, 0x2b, 0x3c},
			wantErr:  false,
		},
		{
			name:     "invalid hex",
			hexStr:   "1a2x3c",
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToBytes(tt.hexStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !compareBytes(got, tt.expected) {
				t.Errorf("HexToBytes() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// compareBytes compares two byte slices for equality.
func compareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
