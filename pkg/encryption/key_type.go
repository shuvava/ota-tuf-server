package encryption

import (
	"encoding/json"
	"strings"

	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// KeyType is a string denoting a public key type
type KeyType string

const (
	// KeyTypeEd25519 is the type of Ed25519 keys.
	KeyTypeEd25519 = KeyType("ed25519")
	// KeyTypeECDSA is the type of ECDSA keys with SHA2 and P256.
	KeyTypeECDSA = KeyType("ecPrime256v1")
	// KeyTypeRSA is the type of RSA keys with RSASSA-PSS and SHA256.
	KeyTypeRSA = KeyType("rsa")
)

// Validate validates if KeyType has correct value
func (k KeyType) Validate() error {
	switch k {
	case KeyTypeEd25519, KeyTypeRSA, KeyTypeECDSA:
		return nil
	default:
		return apperrors.NewAppError(apperrors.ErrorDataValidation, "unsupported key type: "+string(k))
	}
}

// MarshalJSON original key-server return key type in upper case
func (k KeyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToUpper(string(k)))
}

// UnmarshalJSON convert string to lower case
func (k KeyType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	k = ToKeyType(s)
	return nil
}

// ToKeyType converts string to KeyType without validation
func ToKeyType(s string) KeyType {
	s = strings.ToLower(s)
	kt := KeyType(s)
	return kt
}
