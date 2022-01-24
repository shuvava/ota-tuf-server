package encryption

import (
	"github.com/shuvava/go-ota-svc-common/apperrors"
	intData "github.com/shuvava/ota-tuf-server/internal/data"
	"github.com/shuvava/ota-tuf-server/pkg/data"
)

// rawKey is a raw key representation used for marshaling/unmarshaling
type rawKey struct {
	Public  intData.HexBytes `json:"public"`
	Private intData.HexBytes `json:"private,omitempty"`
}

// UnmarshalKey takes key data to a working verifier implementation for the key type.
// This performs any validation over the data.PublicKey to ensure that the verifier is usable
// to verify signatures.
func UnmarshalKey(key *data.Key) (Verifier, error) {
	switch key.Type {
	case data.KeyTypeEd25519:
		return UnmarshalEd25519Key(key)
	case data.KeyTypeRSA:
		return UnmarshalRSAKey(key)
	case data.KeyTypeECDSA:
		return UnmarshalECDSAKey(key)
	}
	return nil, apperrors.NewAppError(apperrors.ErrorDataValidation, "unsupported key type: "+string(key.Type))
}

// NewKey creates a new encryption key of the given type.
func NewKey(keyType data.KeyType) (Key, error) {
	switch keyType {
	case data.KeyTypeEd25519:
		return GenerateEd25519Key()
	case data.KeyTypeRSA:
		return GenerateRSAKey()
	case data.KeyTypeECDSA:
		return GenerateECDSAKey()
	}
	return nil, apperrors.NewAppError(apperrors.ErrorDataValidation, "unsupported key type: "+string(keyType))
}
