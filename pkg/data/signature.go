package data

import (
	"encoding/base64"
	"fmt"

	"github.com/shuvava/ota-tuf-server/pkg/encryption"

	"github.com/shuvava/go-ota-svc-common/apperrors"
	cmndata "github.com/shuvava/go-ota-svc-common/data"
)

type (
	// HashMethod is the type of checksum algorithm
	HashMethod string

	// Checksum hash of some object
	Checksum string

	// Signature is generic signature model
	Signature struct {
		Method encryption.KeyType `json:"method"`
		Value  string             `json:"sig"`
	}

	// ClientSignature is public model with result of object signing
	ClientSignature struct {
		KeyID KeyID `json:"keyid"`
		Signature
	}
)

const (
	// ErrorChecksumValidation checksum validation error
	ErrorChecksumValidation = apperrors.ErrorDataValidation + ":Checksum"
	// ErrorSignatureValidation signature validation error
	ErrorSignatureValidation = apperrors.ErrorDataValidation + ":Signature"
	// ErrorSignatureSerialization data serialization error to JSON
	ErrorSignatureSerialization = apperrors.ErrorDataSerialization + ":Marshal"
)

// NewClientSignature signs data by the key
func NewClientSignature(key encryption.Signer, data []byte) (*ClientSignature, error) {
	sig, err := key.SignMessage(data)
	if err != nil {
		return nil, err
	}

	signature := &ClientSignature{
		KeyID: NewKeyID(key),
	}
	signature.Method = key.Type()
	signature.Value = base64.StdEncoding.EncodeToString(sig)
	return signature, nil
}

// Validate if Checksum has valid format
func (sha Checksum) Validate() error {
	if len(sha) == 0 || !cmndata.ValidHex(64, string(sha)) {
		return apperrors.NewAppError(
			ErrorChecksumValidation,
			fmt.Sprintf("%s must be in hex format 64 charactres long", sha))
	}
	return nil
}

// Validate if Signature has valid format
func (sig *Signature) Validate() error {
	if len(sig.Value) == 0 || !cmndata.ValidBase64(sig.Value) {
		return apperrors.NewAppError(
			ErrorSignatureValidation,
			fmt.Sprintf("%s must be valid base64 string", sig.Value))
	}
	return nil
}

// ToClientSignature converts to ClientSignature
func (sig *Signature) ToClientSignature(keyID KeyID) *ClientSignature {
	csig := &ClientSignature{
		KeyID: keyID,
	}
	csig.Method = sig.Method
	csig.Value = sig.Value
	return csig
}

// ToSignature converts to Signature
func (sig *ClientSignature) ToSignature() *Signature {
	return &Signature{
		Method: sig.Method,
		Value:  sig.Value,
	}
}
