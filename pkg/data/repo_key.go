package data

// RepoKey is a key for a repo
type RepoKey struct {
	// RepoID is the id of the repo
	RepoID RepoID `json:"repo_id"`
	// Role is the role of the key
	Role RoleType `json:"role"`
	// KeyID is the id of the key
	KeyID KeyID `json:"key_id"`
	// Key is the public/private key
	Key Key `json:"key"`
}
