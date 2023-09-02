package encryption_test

import (
	"encoding/json"
	"testing"

	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

func TestECDSAMarshaling(t *testing.T) {
	t.Run("ECDSA key should have correct type", func(t *testing.T) {
		key, _ := encryption.GenerateECDSAKey()
		tkey := key.Type()
		if tkey != encryption.KeyTypeECDSA {
			t.Errorf("key type: is incorrect: %v", tkey)
		}
	})
	t.Run("should be able to unmarshal key", func(t *testing.T) {
		key, _ := encryption.GenerateECDSAKey()
		dtKey, err := key.MarshalAllData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		_, err = encryption.UnmarshalECDSAKey(dtKey)
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
	})
	t.Run("MarshalAllData should work for public key only", func(t *testing.T) {
		key, _ := encryption.GenerateECDSAKey()
		dtKey, err := key.MarshalPublicData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		want, err := encryption.UnmarshalECDSAKey(dtKey)
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		_, err = want.MarshalAllData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
	})
	t.Run("private key should be NOT present in public data", func(t *testing.T) {
		key, _ := encryption.GenerateECDSAKey()
		dtKey, err := key.MarshalPublicData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		want, err := encryption.UnmarshalECDSAKey(dtKey)
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		if want.PrivateKey != nil {
			t.Errorf("priave key should not be in public data")
		}
	})
	t.Run("error should be thrown if bad key", func(t *testing.T) {
		badKeyValue, _ := json.Marshal(true)
		kv := encryption.RawKey{Public: string(badKeyValue)}
		badKey := encryption.SerializedKey{
			Type:  encryption.KeyTypeECDSA,
			Value: kv,
		}
		_, err := encryption.UnmarshalECDSAKey(&badKey)
		if err == nil {
			t.Error("Expected error, but got nil")
		}
	})
}

func TestECDSASignAndVerify(t *testing.T) {
	t.Run("should be able to sign and verify", func(t *testing.T) {
		key, _ := encryption.GenerateECDSAKey()
		message := []byte("hello world")
		signature, err := key.SignMessage(message)
		if err != nil {
			t.Errorf("unable to sign: %v", err)
		}
		if err = key.Verify(message, signature); err != nil {
			t.Errorf("signature is invalid")
		}
	})
	t.Run("should fail if signature does not match key", func(t *testing.T) {
		key, _ := encryption.GenerateECDSAKey()
		keyOther, _ := encryption.GenerateECDSAKey()
		message := []byte("hello world")
		signature, err := key.SignMessage(message)
		if err != nil {
			t.Errorf("unable to sign: %v", err)
		}
		if err = keyOther.Verify(message, signature); err == nil {
			t.Errorf("signature should be invalid")
		}
	})
}
