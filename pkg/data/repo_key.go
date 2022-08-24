package data

import "github.com/shuvava/ota-tuf-server/pkg/encryption"

// RepoKey is a key for a repo
type RepoKey struct {
	// RepoID is the id of the repo
	RepoID RepoID `json:"repo_id"`
	// Role is the role of the key
	Role RoleType `json:"role"`
	// KeyID is the id of the key
	KeyID KeyID `json:"key_id"`
	// Key is the public/private key
	Key encryption.SerializedKey `json:"key"`
}

// ToPublicKey converts internal RepoKey object to the SerializedKey with only public key in it
func (rk RepoKey) ToPublicKey() (*encryption.SerializedKey, error) {
	key, err := rk.Key.UnmarshalKey()
	if err != nil {
		return nil, err
	}
	return key.MarshalPublicData()
}

func (rk RepoKey) ToSinger() (encryption.Signer, error) {
	key, err := rk.Key.UnmarshalKey()
	if err != nil {
		return nil, err
	}
	return key, nil
}
