package txs

import (
	"reflect"
	"testing"
)

func TestTransactionToStringAndFromString(t *testing.T) {
	// Step 1: Initialize a sample Message and Transaction
	content := []byte("Example `personal_sign` message")
	msg := Message{
		Type:    1, // Assuming '1' corresponds to a valid type like ClientRegistrationType
		Content: content,
	}
	transaction := Transaction{
		Msg:       msg,
		Signature: "0x7635254a22cd7762ded328cdb27292884ee2ea21500bbff52ff549b09ef0aaa80dc5e5bc03d67c3ce5176747dcec576bf64d3959bf4a5b0b890fa188026825531b",
	}

	// Step 2: Test ToString method
	encoded, err := transaction.ToString()
	if err != nil {
		t.Fatalf("ToString failed: %v", err)
	}

	// Assert the encoded string is not empty
	if encoded == "" {
		t.Fatal("Encoded string is empty")
	}

	// Step 3: Test FromString method
	var newTransaction Transaction
	if err := newTransaction.FromString(encoded); err != nil {
		t.Fatalf("FromString failed: %v", err)
	}

	// Assert that the decoded transaction matches the original
	if !reflect.DeepEqual(newTransaction, transaction) {
		t.Fatalf("Decoded transaction does not match original. Got %+v, want %+v", newTransaction, transaction)
	}
}
