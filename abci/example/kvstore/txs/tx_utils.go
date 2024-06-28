package txs

import (
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	_ "log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

// HexToBytes converts a hex string to a byte array.
func HexToBytes(hexStr string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
}

// PersonalMessagePrefix is the prefix used by the eth_sign method in MetaMask.
const PersonalMessagePrefix = "\x19Ethereum Signed Message:\n"

// HashPersonalMessage hashes a message according to the personal_sign spec.
func HashPersonalMessage(message []byte) []byte {
	msg := fmt.Sprintf("%s%d%s", PersonalMessagePrefix, len(message), message)
	return crypto.Keccak256([]byte(msg))
}

// RecoverPubKey recovers the public key from a hash and a signature.
func RecoverPubKey(hash, signature []byte) (*ecdsa.PublicKey, error) {
	if len(signature) != 65 {
		return nil, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	if signature[64] != 27 && signature[64] != 28 {
		return nil, fmt.Errorf("invalid signature 'v' value: %d", signature[64])
	}

	// Convert the recovery id (v value) to a format compatible with the secp256k1 library
	signature[64] -= 27

	pubKey, err := secp256k1.RecoverPubkey(hash, signature)
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPubkey(pubKey)
}

// AddressFromPublicKey derives the Ethereum address from a public key.
func AddressFromPublicKey(pubKey *ecdsa.PublicKey) string {
	address := crypto.PubkeyToAddress(*pubKey)
	return address.Hex()
}

func EncodeMessageAndSignature(message string, signature string) (string, error) {
	messageBytes := []byte(message)
	messageLength := uint32(len(messageBytes))

	// Check if the signature has "0x" prefix and remove it if present
	if strings.HasPrefix(signature, "0x") {
		signature = signature[2:]
	}

	// Decode the hex string to a byte slice
	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return "", fmt.Errorf("error decoding signature: %v", err)
	}

	// Create a buffer to hold the length, message, and signature
	buffer := make([]byte, 4+len(messageBytes)+len(signatureBytes))

	// Write the length of the message (4 bytes)
	binary.BigEndian.PutUint32(buffer[:4], messageLength)

	// Write the message bytes
	copy(buffer[4:4+messageLength], messageBytes)

	// Write the signature bytes
	copy(buffer[4+messageLength:], signatureBytes)

	return hex.EncodeToString(buffer), nil
}

func DecodeMessageAndSignature(data string) (string, string, error) {
	// Decode the hex string to a byte slice
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", "", fmt.Errorf("error decoding data: %v", err)
	}

	// Read the message length (4 bytes)
	messageLength := binary.BigEndian.Uint32(dataBytes[:4])

	// Read the message bytes
	messageBytes := dataBytes[4 : 4+messageLength]

	// Read the signature bytes
	signatureBytes := dataBytes[4+messageLength:]

	return string(messageBytes), hex.EncodeToString(signatureBytes), nil
}

/*
func main() {
	// Sample message and signature
	message := "Example `personal_sign` message"
	signature := "0x7635254a22cd776ded328cdb27292884ee2ea21500bbff52ff549b09ef0aaa80dc5e5bc03d"

	// Encode the message and signature
	encoded, err := EncodeMessageAndSignature(message, signature)
	if err != nil {
		fmt.Println("Error encoding message and signature:", err)
		return
	}

	fmt.Println("Encoded String:", encoded)

	// Decode the message and signature
	decodedMessage, decodedSignature, err := DecodeMessageAndSignature(encoded)
	if err != nil {
		fmt.Println("Error decoding message and signature:", err)
		return
	}

	fmt.Println("Decoded Message:", decodedMessage)
	fmt.Println("Decoded Signature:", decodedSignature)
}
*/
