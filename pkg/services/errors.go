package services

import "github.com/shuvava/go-ota-svc-common/apperrors"

const (
	// ErrorSvcSignedContent is SignedContentService group of errors
	ErrorSvcSignedContent = apperrors.ErrorNamespaceSvc + ":SignedContent"
	// ErrorSvcSignedContentKeyNotFound is not found keys error
	ErrorSvcSignedContentKeyNotFound = ErrorSvcSignedContent + ":KeyNotFound"
	// ErrorMissingSignedRole is the error for the case when no signed key found for TuF Repo
	ErrorMissingSignedRole   = ErrorSvcSignedContent + ":MissingSignedRole"
	ErrorVersionAlreadyExist = ErrorSvcSignedContent + ":AlreadyExist"
)
