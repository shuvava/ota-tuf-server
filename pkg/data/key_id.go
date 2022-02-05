package data

import (
	"fmt"

	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/go-ota-svc-common/data"

	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

// KeyID is a type of TUF server key id
// It is internally a fingerprint of the public key
type KeyID string

const ErrorKeyIDValidation = apperrors.ErrorDataValidation + ":KeyID"

func (k KeyID) String() string {
	return string(k)
}

// Validate validates the key id
func (k KeyID) Validate() error {
	kid := k.String()
	if len(kid) == 0 || !data.ValidHex(64, kid) {
		return apperrors.NewAppError(
			ErrorKeyIDValidation,
			fmt.Sprintf("%s must be in hex format 64 charactres long", kid))
	}
	return nil
}

// NewKeyID creates a new KeyID
func NewKeyID(key encryption.Key) KeyID {
	keyID := KeyID(key.FingerprintSHA256())
	return keyID
}

//	KeyIDFromString returns a new KeyID from a string
func KeyIDFromString(s string) (KeyID, error) {
	keyID := KeyID(s)
	if err := keyID.Validate(); err != nil {
		return "", err
	}
	return keyID, nil
}
