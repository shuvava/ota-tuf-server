package errcodes

import "github.com/shuvava/go-ota-svc-common/apperrors"

const (
	// ErrorDataValidationECDSAKey is the error code for ECDSA key validation failure
	ErrorDataValidationECDSAKey = apperrors.ErrorDataValidation + ":ECDSAKey"
	// ErrorDataSerializationECDSAKey is the error code for ECDSA key serialization/creation failure
	ErrorDataSerializationECDSAKey = apperrors.ErrorDataSerialization + ":ECDSAKey"
	// ErrorDataValidationEd25519Key is the error code for Ed25519 key validation failure
	ErrorDataValidationEd25519Key = apperrors.ErrorDataValidation + ":Ed25519Key"
	// ErrorDataSerializationEd25519Key is the error code for Ed25519 key serialization/creation failure
	ErrorDataSerializationEd25519Key = apperrors.ErrorDataSerialization + ":Ed25519Key"
	// ErrorDataValidationRSAKey is the error code for RSA key validation failure
	ErrorDataValidationRSAKey = apperrors.ErrorDataValidation + ":RSAKey"
	// ErrorDataSerializationRSAKey is the error code for RSA key serialization/creation failure
	ErrorDataSerializationRSAKey = apperrors.ErrorDataSerialization + ":RSAKey"
	// ErrorDataSigning is the error code for signing failure
	ErrorDataSigning = apperrors.ErrorNamespaceData + ":Signing"
	// ErrorDataSigningECDSAKey is the error code for ECDSA key signing failure
	ErrorDataSigningECDSAKey = ErrorDataSigning + ":ECDSAKey"
	// ErrorDataSigningEd25519Key is the error code for Ed25519 key signing failure
	ErrorDataSigningEd25519Key = ErrorDataSigning + ":Ed25519Key"
	// ErrorDataSigningRSAKey is the error code for RSA key signing failure
	ErrorDataSigningRSAKey = ErrorDataSigning + ":RSAKey"
)
