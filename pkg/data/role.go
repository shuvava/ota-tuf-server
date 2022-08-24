package data

import (
	"time"

	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

// RoleKeys is sorted list of key associated with RoleType
type RoleKeys struct {
	Keys      []KeyID `json:"keyids"`
	Threshold int     `json:"threshold"`
}

// VersionedRole is a base type of allowed roles
type VersionedRole struct {
	// Version is change number start from 1
	Version uint `json:"version"`
	//ExpiresAt is the expiration time of the role
	ExpiresAt time.Time `json:"expires"`
	RoleType  RoleType  `json:"_type"`
}

// RootRole is VersionedRole with roleType == RoleTypeRoot
type RootRole struct {
	// TODO: maybe it should be just set of KeyID
	// Keys is set of registered keys for the role
	Keys map[KeyID]*encryption.SerializedKey `json:"keys"`
	// KeyRoles is the roles of the keys
	KeyRoles           map[RoleType]RoleKeys `json:"roles"`
	ConsistentSnapshot bool                  `json:"consistent_snapshot"`
	VersionedRole
}

// NewRootRole creates a new RootRole
func NewRootRole(keys map[KeyID]*encryption.SerializedKey, roles map[RoleType][]KeyID, expires time.Time) *RootRole {
	roleKeys := make(map[RoleType]RoleKeys)
	for k, v := range roles {
		roleKeys[k] = RoleKeys{
			Keys:      v,
			Threshold: 1,
		}
	}
	role := &RootRole{
		Keys:               keys,
		KeyRoles:           roleKeys,
		ConsistentSnapshot: false,
	}
	role.Version = 1
	role.ExpiresAt = expires
	role.RoleType = RoleTypeRoot
	return role
}

// GetRoleKeys filters RootRole keys by key role (RoleType)
func (rr *RootRole) GetRoleKeys(rtype RoleType) []KeyID {
	if k, ok := rr.KeyRoles[rtype]; ok {
		return k.Keys
	}
	return nil
}

// Sign create SignedPayload for the RootRole
func (rr *RootRole) Sign(key encryption.Signer) (*SignedPayload[RootRole], error) {
	sig, err := NewClientSignature(key, rr)
	if err != nil {
		return nil, err
	}
	return &SignedPayload[RootRole]{
		Signature: sig,
		Signed:    rr,
	}, nil
}
