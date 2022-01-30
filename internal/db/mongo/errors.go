package mongo

import "github.com/shuvava/go-ota-svc-common/apperrors"

const (
	// ErrorRepoKeyErrorDbAlreadyExist s is the error message for the error when RepoKey is already exist
	ErrorRepoKeyErrorDbAlreadyExist = apperrors.ErrorDbAlreadyExist + ":RepoKey"
)
