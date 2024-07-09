# Transaction System Documentation

This document provides an overview of the transaction system implemented in the provided `message.go` and `transaction.go` files. The system facilitates various types of transactions, each encapsulating specific types of messages with detailed structures and functionalities.

## Table of Contents

1. [Message Types](#message-types)
2. [Transaction Structure](#transaction-structure)
3. [Encoding and Decoding](#encoding-and-decoding)
4. [Examples](#examples)

## Message Types

Messages in the system are designed to handle different scenarios such as client registration, service requests, and miner interactions. Below is a summary of each message type and its structure:

### Defined Message Types

| Message Type             | Value | Description                         |
|--------------------------|-------|-------------------------------------|
| Client Registration      | 1     | Registers a new client              |
| Service Request          | 2     | Request a service by ID             |
| Client Rating            | 3     | Rate a service provided by a miner  |
| Miner Registration       | 4     | Register a new miner                |
| Miner Service Done       | 5     | Notification of completed service   |
| Miner Status Update      | 6     | Update the status of a miner        |
| Miner Reward Claim       | 7     | Claim rewards by the miner          |
| Miner Service Starting   | 8     | Notification of starting a service  |

### Message Structures

#### `ClientRegistrationMsg`

- **ClientName**: The name of the client registering.

#### `ServiceRequestMsg`

- **ServiceID**: Numeric ID of the requested service.
- **Meta**: Metadata related to the service request.

#### `ClientRatingMsg`

- **ReviewedMinerAddr**: Address of the miner being rated.
- **Rating**: Numerical rating given to the miner.

... (Descriptions for other message types would follow a similar format.)

## Transaction Structure

Transactions are composed of a message and a signature, encoded in JSON. The `Transaction` struct encapsulates this information:

### `Transaction` Fields

| Field      | Type     | Description                                 |
|------------|----------|---------------------------------------------|
| Msg        | Message  | Contains the type and content of the message |
| Signature  | string   | Digital signature validating the transaction |

### `Message` Encoding

Messages are base64 encoded to ensure safe transmission over networks. The `Message` struct handles this encoding transparently.

## Encoding and Decoding

Transactions and messages are encoded in JSON format with specific methods for handling base64 content. These methods ensure that messages are securely and accurately transmitted and received.

### Example Methods

- `MarshalJSON()`: Encodes messages with base64 content.
- `UnmarshalJSON()`: Decodes messages from base64 content to their original form.

## Examples

Below are examples of creating and handling transactions within this system:

### Creating a Transaction

```go
msg := ServiceRequestMsg{ServiceID: 1234, Meta: []byte("Details")}
msgBytes, _ := msg.ToBytes()
transaction := Transaction{
    Msg: Message{Type: ServiceRequestType, Content: msgBytes},
    Signature: "exampleSignature",
}
```

### Decoding a Transaction
```go
encodedString := "..."
transaction := Transaction{}
err := transaction.FromString(encodedString)
if err != nil {
    log.Fatal(err)
}
```
