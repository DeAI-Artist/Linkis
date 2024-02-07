package utils

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type Message struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

func TestGetSignersEthAddress_Success(t *testing.T) {
	message := Message{
		Message:   "Hello, this is a test message from DeAI-Artist!",
		Timestamp: time.Now().UTC().Format(time.RFC3339), // Use current time for the timestamp
	}
	jsonString, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	//println(address.String())

	msg := jsonString
	//println(string(jsonString))
	sig, err := GenerateEthSignature(privateKey, msg)
	assert.NoError(t, err)

	recoveredAddress, err := GetSignersEthAddress(msg, sig)
	//println(recoveredAddress.String())
	assert.NoError(t, err)
	assert.Equal(t, address, recoveredAddress)
}

func TestGetSignersEthAddress_InvalidSignatureLength(t *testing.T) {
	msg := []byte("test message")
	sig := []byte("invalid signature length")
	_, err := GetSignersEthAddress(msg, sig)
	assert.EqualError(t, err, "invalid signature: signature length must be 65 bytes")
}

func TestGenerateEthPrefixedMsgHash(t *testing.T) {
	msg := []byte("test message")
	expectedPrefix := "\x19Ethereum Signed Message:\n"
	expectedHash := crypto.Keccak256Hash([]byte(expectedPrefix + "12" + string(msg)))

	hash := GenerateEthPrefixedMsgHash(msg)
	assert.Equal(t, expectedHash, hash)
}

func TestGenerateEthSignature(t *testing.T) {
	message := Message{
		Message:   "Hello, this is a test message from DeAI-Artist!",
		Timestamp: time.Now().UTC().Format(time.RFC3339), // Use current time for the timestamp
	}
	jsonString, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	//println(crypto.PubkeyToAddress(privateKey.PublicKey).String())
	msg := jsonString
	signature, err := GenerateEthSignature(privateKey, msg)
	assert.NoError(t, err)
	assert.Len(t, signature, 65)
	//println(signature)
	recoveredPub, err := crypto.SigToPub(GenerateEthPrefixedMsgHash(msg).Bytes(), signature)
	assert.NoError(t, err)
	assert.Equal(t, privateKey.PublicKey, *recoveredPub)
}
