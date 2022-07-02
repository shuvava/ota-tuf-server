package data

import (
	"github.com/shuvava/go-ota-svc-common/data"
)

// RepoID is a type of TUF server key id
type RepoID data.CorrelationID

// RepoIDNil is a nil value of RepoID
var RepoIDNil = RepoID(data.CorrelationIDNil)

func (r RepoID) String() string {
	return data.CorrelationID(r).String()
}

// MarshalJSON custom json serialization func
func (r RepoID) MarshalJSON() ([]byte, error) {
	c := data.CorrelationID(r)
	return c.MarshalJSON()
}

// UnmarshalJSON custom json deserialization func
func (r RepoID) UnmarshalJSON(d []byte) error {
	c := data.CorrelationID(r)
	return c.UnmarshalJSON(d)
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
