package data

import (
	"time"

	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

// RoleKeys is sorted list of key associated with RoleType
type RoleKeys struct {
	Keys      []KeyID `json:"keyids"`
	Threshold uint    `json:"threshold"`
}

// RepoSigned is a type of signature includes keys and schema version
type RepoSigned struct {
	// TODO: maybe it should be just set of KeyID
	// Keys is set of registered keys for the role
	Keys map[KeyID]*encryption.SerializedKey `json:"keys"`
	// KeyRoles is the roles of the keys
	KeyRoles           map[RoleType]RoleKeys `json:"roles"`
	ConsistentSnapshot bool                  `json:"consistent_snapshot"`
	// Version is change number start from 1
	Version uint `json:"version"`
	//ExpiresAt is the expiration time of the role
	ExpiresAt time.Time `json:"expires"`
	RoleType  RoleType  `json:"_type"`
}

// NewRepoSigned creates a new RepoSigned
func NewRepoSigned(keys map[KeyID]*encryption.SerializedKey, roles map[RoleType][]KeyID, version uint, expires time.Time, threshold uint) *RepoSigned {
	roleKeys := make(map[RoleType]RoleKeys)
	for k, v := range roles {
		roleKeys[k] = RoleKeys{
			Keys:      v,
			Threshold: threshold,
		}
	}
	role := &RepoSigned{
		Keys:               keys,
		KeyRoles:           roleKeys,
		ConsistentSnapshot: false,
	}
	role.Version = version
	role.ExpiresAt = expires
	role.RoleType = RoleTypeRoot
	return role
}

// GetRoleKeys filters RepoSigned keys by key role (RoleType)
func (rr *RepoSigned) GetRoleKeys(rtype RoleType) []KeyID {
	if k, ok := rr.KeyRoles[rtype]; ok {
		return k.Keys
	}
	return nil
}

// Sign create SignedPayload for the RepoSigned
func (rr *RepoSigned) Sign(keys []*RepoKey, threshold uint) (*SignedPayload[RepoSigned], error) {
	signatures, err := SignPayload(rr, keys, threshold)
	if err != nil {
		return nil, err
	}

	return &SignedPayload[RepoSigned]{
		Signatures: signatures,
		Signed:     rr,
	}, nil
}
