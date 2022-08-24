package encryption_test

import (
	"fmt"
	"testing"

	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

func TestUnmarshalPublicKey(t *testing.T) {
	message := []byte("hello world")
	keys := make([]encryption.Key, 3)
	keys[0], _ = encryption.GenerateRSAKey()
	keys[1], _ = encryption.GenerateEd25519Key()
	keys[2], _ = encryption.GenerateECDSAKey()
	for _, key := range keys {
		name := fmt.Sprintf("%s should be able to run UnmarshalKey and verify key", string(key.Type()))
		t.Run(name, func(t *testing.T) {
			skey, err := key.MarshalPublicData()
			if err != nil {
				t.Errorf("unable to marshal key: %v", err)
			}
			signature, err := key.SignMessage(message)
			if err != nil {
				t.Errorf("unable to sign: %v", err)
			}
			v, err := skey.UnmarshalKey()
			if err != nil {
				t.Errorf("unable to unmarshal key: %v", err)
			}
			if err = v.Verify(message, signature); err != nil {
				t.Errorf("signature is invalid")
			}
		})
		name = fmt.Sprintf("%s should should be able to run UnmarshalSigner sign key", string(key.Type()))
		t.Run(name, func(t *testing.T) {
			skey, err := key.MarshalAllData()
			if err != nil {
				t.Errorf("unable to marshal key: %v", err)
			}
			v, err := skey.UnmarshalKey()
			if err != nil {
				t.Errorf("unable to unmarshal key: %v", err)
			}
			signature, err := v.SignMessage(message)
			if err != nil {
				t.Errorf("unable to sign: %v", err)
			}
			if err = key.Verify(message, signature); err != nil {
				t.Errorf("signature is invalid")
			}
		})
	}
}
