package txs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Transaction struct {
	Msg       Message `json:"msg"`
	Signature string  `json:"signature"`
}

type Message struct {
	Type    uint8  `json:"type"`
	Content []byte `json:"content"`
}

// MarshalJSON customizes the JSON encoding of the Message struct
func (m Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	return json.Marshal(&struct {
		Content string `json:"content"`
		*Alias
	}{
		Content: base64.StdEncoding.EncodeToString(m.Content),
		Alias:   (*Alias)(&m),
	})
}

// UnmarshalJSON customizes the JSON decoding of the Message struct
func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		Content string `json:"content"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	var err error
	m.Content, err = base64.StdEncoding.DecodeString(aux.Content)
	return err
}

func (t Transaction) ToString() (string, error) {
	// Marshal the Msg field to JSON
	msgBytes, err := json.Marshal(t.Msg)
	if err != nil {
		return "", fmt.Errorf("error marshaling Msg: %v", err)
	}

	// Convert the JSON bytes to a string
	msgString := string(msgBytes)

	// Use the EncodeMessageAndSignature function to encode the Msg string and the Signature string
	encodedString, err := EncodeMessageAndSignature(msgString, t.Signature)
	if err != nil {
		return "", fmt.Errorf("error encoding message and signature: %v", err)
	}

	return encodedString, nil
}

func (t *Transaction) FromString(encodedString string) error {
	// Use the DecodeMessageAndSignature function to decode the string back to the Msg string and the Signature string
	msgString, signature, err := DecodeMessageAndSignature(encodedString)
	if err != nil {
		return fmt.Errorf("error decoding message and signature: %v", err)
	}

	// Unmarshal the Msg string back to the Message struct
	var msg Message
	if err := json.Unmarshal([]byte(msgString), &msg); err != nil {
		return fmt.Errorf("error unmarshaling Msg: %v", err)
	}

	// Set the decoded values to the Transaction struct
	t.Msg = msg
	t.Signature = signature

	return nil
}

/*
func main() {
	// Create a new message
	originalMessage := Message{
		Type:    1,
		Content: []byte("Hello, World!"),
	}

	// Marshal the message to JSON
	jsonData, err := json.Marshal(originalMessage)
	if err != nil {
		fmt.Println("Error marshaling message:", err)
		return
	}

	fmt.Println("Encoded JSON Data:", string(jsonData))

	// Unmarshal the JSON back to a message
	var decodedMessage Message
	err = json.Unmarshal(jsonData, &decodedMessage)
	if err != nil {
		fmt.Println("Error unmarshaling message:", err)
		return
	}

	fmt.Printf("Decoded Message: %+v\n", decodedMessage)
}
*/
