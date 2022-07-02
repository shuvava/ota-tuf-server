package encryption

import "github.com/shuvava/ota-tuf-server/pkg/data"

// Key represents a common methods of different keys.
type Key interface {
	// Type returns the type of key.
	Type() data.KeyType
	// MarshalAllData returns the data.Key object associated with the verifier contains public and private keys.
	MarshalAllData() (*data.Key, error)
	// MarshalPublicData returns the data.Key object associated with the verifier contains only public key.
	MarshalPublicData() (*data.Key, error)
	// Public this is the public string used as a unique identifier for the verifier instance.
	Public() string
	// FingerprintSHA256 returns the SHA256 fingerprint of the given key.
	FingerprintSHA256() string
}

// NewKeyID creates a new KeyID
func NewKeyID(key Key) data.KeyID {
	keyID := data.KeyID(key.FingerprintSHA256())
	return keyID
}
