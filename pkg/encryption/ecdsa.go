package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"encoding/pem"
	"math/big"
	"strings"

	"github.com/shuvava/go-ota-svc-common/apperrors"

	"github.com/shuvava/ota-tuf-server/pkg/errcodes"
)

type ecdsaSignature struct {
	R, S *big.Int
}

// ECDSAKey is a verifier for ecdsa keys
type ECDSAKey struct {
	Key
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
	keyType    KeyType
}

// NewECDSAKey creates new ECDSAKey
func NewECDSAKey(public *ecdsa.PublicKey, private *ecdsa.PrivateKey) *ECDSAKey {
	return &ECDSAKey{
		PrivateKey: private,
		PublicKey:  public,
		keyType:    KeyTypeECDSA,
	}
}

// GenerateECDSAKey generates a new ecdsa private key and returns it
func GenerateECDSAKey() (*ECDSAKey, error) {
	private, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSerializationECDSAKey, "failed to generate key: ", err)
	}
	signer := NewECDSAKey(&private.PublicKey, private)
	return signer, nil
}

// Type returns the type of key.
func (k *ECDSAKey) Type() KeyType {
	return k.keyType
}

// Method returns the method of signature
func (k *ECDSAKey) Method() KeyMethod {
	return KeyMethodECPrime256V1
}

// MarshalAllData returns the data.SerializedKey object associated with the verifier contains public and private keys.
func (k *ECDSAKey) MarshalAllData() (*SerializedKey, error) {
	key := RawKey{}
	if k.PrivateKey != nil {
		pri := k.PrivateKey.D.Bytes()
		pkey := encodePrivateKey(pri)
		key.Private = &pkey
	}

	return k.marshalKey(key)
}

// MarshalPublicData returns the data.SerializedKey object associated with the verifier contains only public key.
func (k *ECDSAKey) MarshalPublicData() (*SerializedKey, error) {
	return k.marshalKey(RawKey{})
}

// Public this is the public string used as a unique identifier for the verifier instance.
func (k *ECDSAKey) Public() string {
	pub := elliptic.MarshalCompressed(k.PublicKey.Curve, k.PublicKey.X, k.PublicKey.Y)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pub,
	})
	return string(pubBytes)
}

// SignMessage signs a message with the private key.
func (k *ECDSAKey) SignMessage(message []byte) ([]byte, error) {
	hash := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, k.PrivateKey, hash[:])
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSigningECDSAKey, "failed to sign message: ", err)
	}
	sigToMarshal := ecdsaSignature{R: r, S: s}
	keySig, err := asn1.Marshal(sigToMarshal)
	if err != nil {
		return nil, apperrors.CreateError(errcodes.ErrorDataSigningECDSAKey, "failed to sign message: ", err)
	}
	return keySig, nil
}

// Verify takes a message and signature, all as byte slices,
// and determines whether the signature is valid for the given
// key and message.
func (k *ECDSAKey) Verify(msg, sig []byte) error {
	var signature ecdsaSignature
	if _, err := asn1.Unmarshal(sig, &signature); err != nil {
		return err
	}

	hash := sha256.Sum256(msg)
	if !ecdsa.Verify(k.PublicKey, hash[:], signature.R, signature.S) {
		return apperrors.NewAppError(errcodes.ErrorDataValidationECDSAKey, "tuf: ecdsa signature verification failed")
	}
	return nil
}

// VerifyECDSAKey is a helper function to verify an ecdsa key.
func VerifyECDSAKey(v *ECDSAKey) error {
	if _, err := v.PublicKey.ECDH(); err != nil {
		return apperrors.NewAppError(errcodes.ErrorDataValidationECDSAKey, "tuf: ecdsa key is invalid")
	}
	return nil
}

// UnmarshalECDSAKey is a helper function to unmarshal an ecdsa key from a data.SerializedKey.
func UnmarshalECDSAKey(key *SerializedKey) (*ECDSAKey, error) {
	kv := key.Value
	pub, err := decodePublicKey(kv.Public)
	if err != nil {
		return nil, apperrors.NewAppError(errcodes.ErrorDataValidationECDSAKey, "Unable to decode PEM block in public key")
	}
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), pub)
	if x == nil || y == nil {
		return nil, apperrors.NewAppError(errcodes.ErrorDataSerializationECDSAKey, "error on unmarshalling key")
	}
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	ecdsaKey := NewECDSAKey(&publicKey, nil)
	if kv.Private != nil {
		if pri, err := decodePrivateKey(kv.Private); err == nil {
			privateKey := ecdsa.PrivateKey{
				PublicKey: publicKey,
				D:         new(big.Int).SetBytes(pri),
			}
			ecdsaKey.PrivateKey = &privateKey
		}

	}

	if err := VerifyECDSAKey(ecdsaKey); err != nil {
		return nil, err
	}
	return ecdsaKey, nil
}

func (k *ECDSAKey) marshalKey(kv RawKey) (*SerializedKey, error) {
	kv.Public = k.Public()

	return &SerializedKey{
		Type:  k.keyType,
		Value: kv,
	}, nil
}

// FingerprintSHA256 returns the SHA256 hex fingerprint of the public key.
func (k *ECDSAKey) FingerprintSHA256() string {
	hash := sha256.Sum256(elliptic.MarshalCompressed(k.PublicKey.Curve, k.PublicKey.X, k.PublicKey.Y))
	return strings.ToLower(hex.EncodeToString(hash[:]))
}
