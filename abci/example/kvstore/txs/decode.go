package txs

import (
	"crypto/ecdsa"
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
