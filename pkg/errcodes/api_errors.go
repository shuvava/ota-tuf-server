package errcodes

import "github.com/shuvava/go-ota-svc-common/apperrors"

const (
	errorAPI = apperrors.ErrorDataValidation + ":API"
	// ErrorAPIRequestValidationParamMissing is the error code for validation API for missing required param
	ErrorAPIRequestValidationParamMissing = errorAPI + ":ParamMissing"
	// ErrorAPIRequestValidationParamInvalid is the error code for validation API for invalid param type
	ErrorAPIRequestValidationParamInvalid = errorAPI + ":ParamInvalid"
	// ErrorAPIRequestValidationBodyValidation is the error code for validation API for invalid body
	ErrorAPIRequestValidationBodyValidation = errorAPI + ":BodyInvalid"
)
