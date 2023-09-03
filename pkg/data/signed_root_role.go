package data

import (
	"time"
)

// SignedRootRole is a TUF repo signed Role
type SignedRootRole struct {
	// RepoID is the id of the repo
	RepoID    RepoID                   `json:"repo_id"`
	ExpiresAt time.Time                `json:"expires_at"`
	Version   uint                     `json:"version"`
	Threshold uint                     `json:"threshold"`
	Content   *SignedPayload[RootRole] `json:"signed_payload"`
}

// ShouldBeRenewed TUF repo signature should be renewed
func (sig *SignedRootRole) ShouldBeRenewed() bool {
	return sig.ExpiresAt.Before(time.Now().Add(time.Hour).UTC())
}
