package data

import (
	"fmt"

	"github.com/shuvava/go-ota-svc-common/data"
)

// KeyID is a type of TUF server key id
type KeyID string

//Validate if Commit has valid format
func (key KeyID) Validate() error {
	err := fmt.Errorf("%s is not a sha-256 commit hash", key)
	sha := string(key)
	if !data.ValidHex(64, sha) {
		return err
	}

	return nil
}

// NewKeyID creates a new KeyID
func NewKeyID(id string) (KeyID, error) {
	key := KeyID(id)
	if err := key.Validate(); err != nil {
		return "", err
	}
	return key, nil
}
