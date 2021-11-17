package data

import (
	"encoding/json"
)

// Key is common struct for signature and encryption keys
type Key struct {
	// Type is key type
	Type KeyType `json:"keytype"`
	// Value is key value
	Value json.RawMessage `json:"keyval"`
}

// PrivateKey is a private key
type PrivateKey struct {
	Key
}

// PublicKey is a public key
type PublicKey struct {
	Key
	ids []string
}
