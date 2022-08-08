package encryption_test

import (
	"testing"

	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

func TestUnmarshalPublicKey(t *testing.T) {
	message := []byte("hello world")
	t.Run("RSA should be able to UnmarshalPublicKey and verify key", func(t *testing.T) {
		key, _ := encryption.GenerateRSAKey()
		skey, err := key.MarshalPublicData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		signature, err := key.SignMessage(message)
		if err != nil {
			t.Errorf("unable to sign: %v", err)
		}
		v, err := skey.UnmarshalPublicKey()
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		if err = v.Verify(message, signature); err != nil {
			t.Errorf("signature is invalid")
		}
	})
	t.Run("Ed25519 should be able to UnmarshalPublicKey and verify key", func(t *testing.T) {
		key, _ := encryption.GenerateEd25519Key()
		skey, err := key.MarshalPublicData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		signature, err := key.SignMessage(message)
		if err != nil {
			t.Errorf("unable to sign: %v", err)
		}
		v, err := skey.UnmarshalPublicKey()
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		if err = v.Verify(message, signature); err != nil {
			t.Errorf("signature is invalid")
		}
	})
	t.Run("ECDSA should be able to UnmarshalPublicKey and verify key", func(t *testing.T) {
		key, _ := encryption.GenerateECDSAKey()
		skey, err := key.MarshalPublicData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		signature, err := key.SignMessage(message)
		if err != nil {
			t.Errorf("unable to sign: %v", err)
		}
		v, err := skey.UnmarshalPublicKey()
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		if err = v.Verify(message, signature); err != nil {
			t.Errorf("signature is invalid")
		}
	})
}
