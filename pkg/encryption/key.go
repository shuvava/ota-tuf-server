package encryption

import "github.com/shuvava/ota-tuf-server/pkg/data"

// Key represents a common methods of different keys.
type Key interface {
	// Type returns the type of key.
	Type() data.KeyType
	// MarshalAllData returns the data.Key object associated with the verifier contains public and private keys.
	MarshalAllData() (*data.Key, error)
}
