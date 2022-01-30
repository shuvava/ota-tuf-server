package data

import (
	"github.com/shuvava/go-ota-svc-common/data"
)

// KeyID is a type of TUF server key id
type KeyID data.CorrelationID

func (key KeyID) String() string {
	return data.CorrelationID(key).String()
}

// NewKeyID creates a new KeyID
func NewKeyID(repoID RepoID, roleType RoleType) KeyID {
	id := data.NewChildCorrelationID(data.CorrelationID(repoID), string(roleType))
	key := KeyID(id)
	return key
}

//	KeyIDFromString returns a new KeyID from a string
func KeyIDFromString(s string) (KeyID, error) {
	id, err := data.CorrelationIDFromString(s)
	if err != nil {
		return KeyID(id), err
	}
	return KeyID(id), nil
}
