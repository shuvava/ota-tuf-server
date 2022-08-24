package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"strings"

	"github.com/shuvava/go-ota-svc-common/apperrors"

	"github.com/shuvava/ota-tuf-server/pkg/errcodes"
)

// RSAKey is
type RSAKey struct {
	Key
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
	keyType    KeyType
}

// NewRSAKey creates new RSAKey
func NewRSAKey(public *rsa.PublicKey, private *rsa.PrivateKey) *RSAKey {
	return &RSAKey{
		PrivateKey: private,
		PublicKey:  public,
		keyType:    KeyTypeRSA,
	}
}

// GenerateRSAKey generates a new rsa private key and returns it
func GenerateRSAKey() (*RSAKey, error) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSerializationRSAKey, "failed to generate key: ", err)
	}
	key := NewRSAKey(private.Public().(*rsa.PublicKey), private)
	return key, nil
}

// Type returns the type of key.
func (k *RSAKey) Type() KeyType {
	return k.keyType
}

// MarshalPublicData returns the data.SerializedKey object associated with the verifier contains only public key.
func (k *RSAKey) MarshalPublicData() (*SerializedKey, error) {
	return k.marshalKey(rawKey{})
}

// MarshalAllData returns the data.SerializedKey object associated with the verifier contains public and private keys.
func (k *RSAKey) MarshalAllData() (*SerializedKey, error) {
	var priBytes []byte
	if k.PrivateKey != nil {
		pri := x509.MarshalPKCS1PrivateKey(k.PrivateKey)
		priBytes = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pri,
		})
	}

	return k.marshalKey(rawKey{
		Private: priBytes,
	})
}

// Public this is the public string used as a unique identifier for the verifier instance.
func (k *RSAKey) Public() string {
	pub, _ := x509.MarshalPKIXPublicKey(k.PublicKey)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pub,
	})
	return string(pubBytes)
}

// SignMessage signs a message with the private key.
func (k *RSAKey) SignMessage(message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)
	keySig, err := rsa.SignPSS(rand.Reader, k.PrivateKey, crypto.SHA256, hash[:], &rsa.PSSOptions{})
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSerializationRSAKey, "failed to sign message: ", err)
	}
	return keySig, nil
}

// Verify takes a message and signature, all as byte slices,
// and determines whether the signature is valid for the given
// key and message.
func (k *RSAKey) Verify(msg, sig []byte) error {
	hash := sha256.Sum256(msg)

	if err := rsa.VerifyPSS(k.PublicKey, crypto.SHA256, hash[:], sig, &rsa.PSSOptions{}); err != nil {
		return apperrors.CreateError(errcodes.ErrorDataSerializationRSAKey, "failed to verify signature: ", err)
	}
	return nil
}

// UnmarshalRSAKey is a helper function to unmarshal an RSA key from a data.SerializedKey.
func UnmarshalRSAKey(key *SerializedKey) (*RSAKey, error) {
	var kv rawKey
	if err := json.Unmarshal(key.Value, &kv); err != nil {
		return nil, err
	}
	block, _ := pem.Decode(kv.Public)
	if block == nil {
		return nil, apperrors.NewAppError(errcodes.ErrorDataValidationRSAKey, "Unable to decode PEM block in public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataValidationRSAKey, "failed to unmarshal public key: ", err)
	}
	rsaKey := NewRSAKey(publicKey.(*rsa.PublicKey), nil)

	block, _ = pem.Decode(kv.Private)
	if block != nil {
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, apperrors.CreateError(errcodes.ErrorDataValidationRSAKey, "failed to unmarshal private key: ", err)
		}
		rsaKey.PrivateKey = privateKey
	}

	if err := VerifyRSAKey(rsaKey); err != nil {
		return nil, err
	}
	return rsaKey, nil
}

// VerifyRSAKey is a helper function to verify a rsa key.
func VerifyRSAKey(v *RSAKey) error {
	if v.PublicKey == nil {
		return apperrors.NewAppError(errcodes.ErrorDataValidationRSAKey, "public key is nil")
	}
	return nil
}

func (k *RSAKey) marshalKey(kv rawKey) (*SerializedKey, error) {
	pub, err := x509.MarshalPKIXPublicKey(k.PublicKey)
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataValidationRSAKey, "failed to marshal public key: ", err)
	}
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pub,
	})
	kv.Public = pubBytes

	valueBytes, err := json.Marshal(kv)
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataValidationRSAKey, "failed to marshal key: ", err)
	}

	return &SerializedKey{
		Type:  k.keyType,
		Value: valueBytes,
	}, nil
}

// FingerprintSHA256 returns the SHA256 hex fingerprint of the public key.
func (k *RSAKey) FingerprintSHA256() string {
	hash := sha256.Sum256(k.PublicKey.N.Bytes())
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
