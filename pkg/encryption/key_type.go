package encryption

import "github.com/shuvava/go-ota-svc-common/apperrors"

// KeyType is a string denoting a public key signature system
type KeyType string

const (
	// KeyTypeEd25519 is the type of Ed25519 keys.
	KeyTypeEd25519 = KeyType("ed25519")
	// KeyTypeECDSA is the type of ECDSA keys with SHA2 and P256.
	KeyTypeECDSA = KeyType("ecPrime256v1")
	// KeyTypeRSA is the type of RSA keys with RSASSA-PSS and SHA256.
	KeyTypeRSA = KeyType("rsassa-pss-sha256")
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
