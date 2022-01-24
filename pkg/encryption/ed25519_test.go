package encryption_test

import (
	"encoding/json"
	"testing"

	intData "github.com/shuvava/ota-tuf-server/internal/data"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

func TestEd25519Marshaling(t *testing.T) {
	t.Run("should be able to unmarshal key", func(t *testing.T) {
		key, _ := encryption.GenerateEd25519Key()
		dtKey, err := key.MarshalAllData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		want, err := encryption.UnmarshalEd25519Key(dtKey)
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		if intData.HexBytes(want.PrivateKey).String() != intData.HexBytes(key.PrivateKey).String() {
			t.Errorf("want: %v, got: %v", key.PrivateKey, want.PrivateKey)
		}
	})
	t.Run("private key should be NOT present in public data", func(t *testing.T) {
		key, _ := encryption.GenerateEd25519Key()
		dtKey, err := key.MarshalPublicData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		want, err := encryption.UnmarshalEd25519Key(dtKey)
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		if len(want.PrivateKey) > 0 {
			t.Errorf("priave key should not be in public data")
		}
	})
	t.Run("error should be thrown if bad key", func(t *testing.T) {
		badKeyValue, _ := json.Marshal(true)
		badKey := data.Key{
			Type:  data.KeyTypeEd25519,
			Value: badKeyValue,
		}
		_, err := encryption.UnmarshalEd25519Key(&badKey)
		if err == nil {
			t.Error("Expected error, but got nil")
		}
	})
}

func TestEd25519SignAndVerify(t *testing.T) {
	t.Run("should be able to sign and verify", func(t *testing.T) {
		key, _ := encryption.GenerateEd25519Key()
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
		key, _ := encryption.GenerateEd25519Key()
		keyOther, _ := encryption.GenerateEd25519Key()
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
