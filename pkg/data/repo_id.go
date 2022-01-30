package data

import "github.com/shuvava/go-ota-svc-common/data"

// RepoID is a type of TUF server key id
type RepoID data.CorrelationID

func (r RepoID) String() string {
	return data.CorrelationID(r).String()
}

// NewRepoID returns a new RepoID
func NewRepoID() RepoID {
	return RepoID(data.NewCorrelationID())
}

// RepoIDFromString returns a new RepoID from a string
func RepoIDFromString(s string) (RepoID, error) {
	id, err := data.CorrelationIDFromString(s)
	return RepoID(id), err
}
