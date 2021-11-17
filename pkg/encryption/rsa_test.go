package encryption_test

import (
	"encoding/json"
	"testing"

	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

func TestRSAMarshaling(t *testing.T) {
	t.Run("should be able to unmarshal key", func(t *testing.T) {
		key, _ := encryption.GenerateRSAKey()
		dtKey, err := key.MarshalAllData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		_, err = encryption.UnmarshalRSAKey(dtKey)
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
	})
	t.Run("private key should be NOT present in public data", func(t *testing.T) {
		key, _ := encryption.GenerateRSAKey()
		dtKey, err := key.MarshalPublicData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		want, err := encryption.UnmarshalRSAKey(dtKey)
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		if want.PrivateKey != nil {
			t.Errorf("priave key should not be in public data")
		}
	})
	t.Run("error should be thrown if bad key", func(t *testing.T) {
		badKeyValue, _ := json.Marshal(true)
		badKey := data.Key{
			Type:  data.KeyTypeRSA,
			Value: badKeyValue,
		}
		_, err := encryption.UnmarshalRSAKey(&badKey)
		if err == nil {
			t.Error("Expected error, but got nil")
		}
	})
}

func TestRSASignAndVerify(t *testing.T) {
	t.Run("should be able to sign and verify", func(t *testing.T) {
		key, _ := encryption.GenerateRSAKey()
		message := []byte("hello world")
		signature, err := key.SignMessage(message)
		if err != nil {
			t.Errorf("unable to sign: %v", err)
		}
		if err := key.Verify(message, signature); err != nil {
			t.Errorf("signature is invalid")
		}
	})
}
