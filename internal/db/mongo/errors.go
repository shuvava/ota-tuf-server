package mongo

import "github.com/shuvava/go-ota-svc-common/apperrors"

const (
	// ErrorSignedContentErrorDbAlreadyExist is the error message for the error when data.SignedRootRole is already exist
	ErrorSignedContentErrorDbAlreadyExist = apperrors.ErrorDbAlreadyExist + ":SignedContent"
	// ErrorRepoKeyErrorDbAlreadyExist is the error message for the error when data.RepoKey is already exist
	ErrorRepoKeyErrorDbAlreadyExist = apperrors.ErrorDbAlreadyExist + ":RepoKey"
	// ErrorRepoErrorDbAlreadyExist is the error message for the error when data.Repo is already exist
	ErrorRepoErrorDbAlreadyExist = apperrors.ErrorDbAlreadyExist + ":Repo"
)
