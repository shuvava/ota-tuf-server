package errcodes

import "github.com/shuvava/go-ota-svc-common/apperrors"

const (
	// ErrorDataValidationECDSAKey is the error code for ECDSA key validation failure
	ErrorDataValidationECDSAKey = apperrors.ErrorDataValidation + ":ECDSAKey"
)
