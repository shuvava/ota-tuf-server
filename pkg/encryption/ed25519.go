package encryption

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"

	"github.com/shuvava/go-ota-svc-common/apperrors"

	intData "github.com/shuvava/ota-tuf-server/internal/data"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/errcodes"
)

// Ed25519Key is a verifier for ed25519 keys
type Ed25519Key struct {
	Key
	Verifier
	Signer
	ed25519.PrivateKey
	PublicKey ed25519.PublicKey
	keyType   data.KeyType
}

// GenerateEd25519Key generates a new ed25519 private key and returns it
func GenerateEd25519Key() (*Ed25519Key, error) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSerializationEd25519Key, "failed to generate key: ", err)
	}
	signer := Ed25519Key{
		PrivateKey: private,
		PublicKey:  public,
		keyType:    data.KeyTypeEd25519,
	}
	return &signer, nil
}

// Type returns the type of the signature scheme.
func (k *Ed25519Key) Type() data.KeyType {
	return k.keyType
}

// MarshalPublicData returns the data.PublicKey object associated with the verifier.
func (k *Ed25519Key) MarshalPublicData() (*data.Key, error) {
	return k.marshalKey(rawKey{})
}

// MarshalAllData returns the data.Key object associated with the verifier contains public and private keys.
func (k *Ed25519Key) MarshalAllData() (*data.Key, error) {
	kv := rawKey{
		Private: intData.HexBytes(k.PrivateKey),
	}

	return k.marshalKey(kv)
}

// Public this is the public string used as a unique identifier for the verifier instance.
func (k *Ed25519Key) Public() string {
	return string(k.PublicKey)
}

// Verify takes a message and signature, all as byte slices,
// and determines whether the signature is valid for the given
// key and message.
func (k *Ed25519Key) Verify(msg, sig []byte) error {
	if !ed25519.Verify(k.PublicKey, msg, sig) {
		return apperrors.NewAppError(apperrors.ErrorDataValidation, "tuf: ed25519 signature verification failed")
	}
	return nil
}

// SignMessage signs a message with the private key.
func (k *Ed25519Key) SignMessage(message []byte) ([]byte, error) {
	keySig, err := k.Sign(rand.Reader, message, crypto.Hash(0))
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSigningEd25519Key, "failed to sign message: ", err)
	}
	return keySig, nil
}

// UnmarshalEd25519Key is a helper function to unmarshal an ed25519 key from a data.Key.
func UnmarshalEd25519Key(key *data.Key) (*Ed25519Key, error) {
	var kv rawKey
	if err := json.Unmarshal(key.Value, &kv); err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSerializationEd25519Key, "failed to deserialize key: ", err)
	}

	privateKey := ed25519.PrivateKey(kv.Private)
	publicKey := ed25519.PublicKey(kv.Public)
	verifier := Ed25519Key{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		keyType:    data.KeyTypeEd25519,
	}

	if err := VerifyEd25519Key(&verifier); err != nil {
		return nil, err
	}
	return &verifier, nil
}

// VerifyEd25519Key is a helper function to verify an ed25519 key.
func VerifyEd25519Key(v *Ed25519Key) error {
	if len(v.PublicKey) != ed25519.PublicKeySize {
		return apperrors.NewAppError(errcodes.ErrorDataValidationEd25519Key, "tuf: ed25519 public key is invalid")
	}
	if v.PrivateKey != nil && len(v.PrivateKey) != ed25519.PrivateKeySize {
		return apperrors.NewAppError(errcodes.ErrorDataValidationEd25519Key, "tuf: ed25519 private key is invalid")
	}
	if v.PrivateKey != nil && !v.PublicKey.Equal(v.PrivateKey.Public().(ed25519.PublicKey)) {
		return apperrors.NewAppError(errcodes.ErrorDataValidationEd25519Key, "tuf: ed25519 public key does not match private key")
	}
	return nil
}

func (k *Ed25519Key) marshalKey(kv rawKey) (*data.Key, error) {
	kv.Public = intData.HexBytes(k.PublicKey)

	valueBytes, err := json.Marshal(kv)
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSerializationEd25519Key, "failed to marshal key: ", err)
	}

	return &data.Key{
		Type:  k.keyType,
		Value: valueBytes,
	}, nil
}
