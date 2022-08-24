package encryption

import (
	"encoding/json"

	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// SerializedKey is common struct for signature and encryption keys
type SerializedKey struct {
	// Type is key type
	Type KeyType `json:"keytype"`
	// Value is key value
	Value json.RawMessage `json:"keyval"`
}

// UnmarshalKey takes key data and convert it to valid Key type
// Node if SerializedKey includes only PublicKey only Verify interface will work
func (key *SerializedKey) UnmarshalKey() (Key, error) {
	switch key.Type {
	case KeyTypeEd25519:
		return UnmarshalEd25519Key(key)
	case KeyTypeRSA:
		return UnmarshalRSAKey(key)
	case KeyTypeECDSA:
		return UnmarshalECDSAKey(key)
	}
	return nil, apperrors.NewAppError(apperrors.ErrorDataValidation, "unsupported key type: "+string(key.Type))
}

/*
// PrivateKey is a private key
type PrivateKey struct {
	SerializedKey
}

// PublicKey is a public key
type PublicKey struct {
	SerializedKey
	ids []string
}
*/
