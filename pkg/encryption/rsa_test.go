package encryption_test

import (
	"encoding/json"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
	"testing"
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
	t.Run("MarshalAllData should work for public key only", func(t *testing.T) {
		key, _ := encryption.GenerateRSAKey()
		dtKey, err := key.MarshalPublicData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
		}
		want, err := encryption.UnmarshalRSAKey(dtKey)
		if err != nil {
			t.Errorf("unable to unmarshal key: %v", err)
		}
		_, err = want.MarshalAllData()
		if err != nil {
			t.Errorf("unable to marshal key: %v", err)
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
		if err = key.Verify(message, signature); err != nil {
			t.Errorf("signature is invalid")
		}
	})
	t.Run("should fail if signature does not match key", func(t *testing.T) {
		key, _ := encryption.GenerateRSAKey()
		keyOther, _ := encryption.GenerateRSAKey()
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

func TestRSAKey_FingerprintSHA256(t *testing.T) {
	t.Run("should be able to generate fingerprint", func(t *testing.T) {
		key, _ := encryption.GenerateRSAKey()
		fingerprint := key.FingerprintSHA256()
		if fingerprint == "" {
			t.Errorf("fingerprint is empty")
		}
		if len(fingerprint) != 64 {
			t.Errorf("fingerprint is not 64 characters")
		}
	})
	//t.Run("should generate valid fingerprint", func(t *testing.T) {
	//	key, _ := encryption.GenerateRSAKey()
	//	// Generate a pem block with the private key
	//	keyPem := pem.EncodeToMemory(&pem.Block{
	//		Type:  "RSA PRIVATE KEY",
	//		Bytes: x509.MarshalPKCS1PrivateKey(key.PrivateKey),
	//	})
	//	tml := x509.Certificate{
	//		// you can add any attr that you need
	//		NotBefore: time.Now(),
	//		NotAfter:  time.Now().AddDate(5, 0, 0),
	//		// you have to generate a different serial number each execution
	//		SerialNumber: big.NewInt(123123),
	//		Subject: pkix.Name{
	//			CommonName:   "New Name",
	//			Organization: []string{"New Org."},
	//		},
	//		BasicConstraintsValid: true,
	//	}
	//	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)
	//	if err != nil {
	//		t.Errorf("Certificate cannot be created.")
	//	}
	//})
}
