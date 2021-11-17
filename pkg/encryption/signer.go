package encryption

import (
	"github.com/shuvava/ota-tuf-server/pkg/data"
)

// Signer is an interface for an opaque private key that can be used for signing operations.
type Signer interface {
	// MarshalAllData returns the data.Key object associated with the verifier contains public and private keys.
	MarshalAllData() (*data.Key, error)
	// SignMessage signs a message with the private key.
	SignMessage(message []byte) ([]byte, error)
}
