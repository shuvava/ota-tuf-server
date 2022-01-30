package data

import (
	"time"

	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// RoleType is top-level roles in the framework
type RoleType string

const (
	// RoleTypeRoot is The root role delegates trust to specific keys trusted for all other top-level roles used in the system.
	// https://theupdateframework.github.io/specification/latest/#root
	RoleTypeRoot = RoleType("root")
	// RoleTypeTargets is The targets role’s signature indicates which target files are trusted by clients. The targets role signs metadata that describes these files, not the actual target files themselves.
	// https://theupdateframework.github.io/specification/latest/#targets
	RoleTypeTargets = RoleType("targets")
	// RoleTypeSnapshot is The snapshot role signs a metadata file that provides information about the latest version of all targets metadata on the repository
	// https://theupdateframework.github.io/specification/latest/#snapshot
	RoleTypeSnapshot = RoleType("snapshot")
	// RoleTypeTimestamp is To prevent an adversary from replaying an out-of-date signed metadata file whose signature has not yet expired, an automated process periodically signs a timestamped statement containing the hash of the snapshot file.
	// https://theupdateframework.github.io/specification/latest/#timestamp
	RoleTypeTimestamp = RoleType("timestamp")
)

// TopLevelRoles is a list of top-level roles defined in the specification
// https://theupdateframework.github.io/specification/latest/#roles-and-pki
var TopLevelRoles = map[RoleType]struct{}{
	RoleTypeRoot:      {},
	RoleTypeTargets:   {},
	RoleTypeSnapshot:  {},
	RoleTypeTimestamp: {},
}

// DefaultExpires returns the default expiration time for a role
func DefaultExpires(role RoleType) time.Time {
	var t time.Time
	switch role {
	case RoleTypeRoot:
		t = time.Now().AddDate(1, 0, 0)
	case RoleTypeTargets:
		t = time.Now().AddDate(0, 3, 0)
	case RoleTypeSnapshot:
		t = time.Now().AddDate(0, 0, 7)
	case RoleTypeTimestamp:
		t = time.Now().AddDate(0, 0, 1)
	}
	return t.UTC().Round(time.Second)
}

// NewRoleType returns a new RoleType from a string
func NewRoleType(name string) (RoleType, error) {
	role := RoleType(name)
	if _, ok := TopLevelRoles[role]; !ok {
		return "", apperrors.NewAppError(apperrors.ErrorDataValidation, "tuf: invalid role '"+name+"'")
	}
	return role, nil
}
