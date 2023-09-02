package encryption

import (
	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// NewKey creates a new encryption key of the given type.
func NewKey(keyType KeyType) (Key, error) {
	switch keyType {
	case KeyTypeEd25519:
		return GenerateEd25519Key()
	case KeyTypeRSA:
		return GenerateRSAKey()
	case KeyTypeECDSA:
		return GenerateECDSAKey()
	}
	return nil, apperrors.NewAppError(apperrors.ErrorDataValidation, "unsupported key type: "+string(keyType))
}

// FingerprintSHA256 returns the SHA256 fingerprint of the given key.
func FingerprintSHA256(key Key) string {
	switch key.Type() {
	case KeyTypeEd25519:
		k := key.(*Ed25519Key)
		return k.FingerprintSHA256()
	case KeyTypeRSA:
		k := key.(*RSAKey)
		return k.FingerprintSHA256()
	case KeyTypeECDSA:
		k := key.(*ECDSAKey)
		return k.FingerprintSHA256()
	}
	return ""
}
