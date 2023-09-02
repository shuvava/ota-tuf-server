package encryption

import (
	"encoding/pem"

	"github.com/shuvava/go-ota-svc-common/apperrors"
)

const (
	// PEMTypePublicKey is the public key pem block type
	PEMTypePublicKey = "PUBLIC KEY"
	// PEMTypePrivateKey is the private key pem block type
	PEMTypePrivateKey = "PRIVATE KEY"
)

// RawKey is a raw key representation used for marshaling/unmarshaling
type RawKey struct {
	Public  string  `json:"public"`
	Private *string `json:"private,omitempty"`
}

// SerializedKey is common struct for signature and encryption keys
type SerializedKey struct {
	// Type is key type
	Type KeyType `json:"keytype"`
	// Value is key value
	Value RawKey `json:"keyval"`
}

// UnmarshalKey takes key data and convert it to valid Key type
// Node if SerializedKey includes only PublicKey only Verify interface will work
func (key *SerializedKey) UnmarshalKey() (Key, error) {
	switch key.Type {
	case KeyTypeEd25519:
		return UnmarshalEd25519Key(key)
	case KeyTypeRSA:
		return UnmarshalRSAKey(key)
	case KeyTypeECDSA:
		return UnmarshalECDSAKey(key)
	}
	return nil, apperrors.NewAppError(apperrors.ErrorDataValidation, "unsupported key type: "+string(key.Type))
}

func encodePublicKey(k []byte) string {
	return encodeKey(k, PEMTypePublicKey)
}

func encodePrivateKey(k []byte) string {
	return encodeKey(k, PEMTypePrivateKey)
}

func decodePublicKey(k string) ([]byte, error) {
	return decodeKey(k, PEMTypePublicKey)
}

func decodePrivateKey(k *string) ([]byte, error) {
	return decodeKey(*k, PEMTypePrivateKey)
}

func encodeKey(k []byte, ktype string) string {
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  ktype,
		Bytes: k,
	})
	return string(pubBytes)
}

func decodeKey(k string, ktype string) ([]byte, error) {
	pemBlock := []byte(k)
	var derBlock *pem.Block
	for {
		derBlock, pemBlock = pem.Decode(pemBlock)
		if derBlock == nil {
			return nil, apperrors.NewAppError(apperrors.ErrorDataSerialization, "Unable to decode PEM block in public key")
		}
		if derBlock.Type == ktype {
			return derBlock.Bytes, nil
		}
	}
}
