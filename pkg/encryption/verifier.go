package encryption

import (
	"github.com/shuvava/ota-tuf-server/pkg/data"
)

// A Verifier verifies public key signatures.
type Verifier interface {
	// MarshalPublicData returns the data.Key object associated with the verifier contains only public key.
	MarshalPublicData() (*data.Key, error)
	// Public this is the public string used as a unique identifier for the verifier instance.
	Public() string
	// Verify takes a message and signature, all as byte slices,
	// and determines whether the signature is valid for the given
	// key and message.
	Verify(msg, sig []byte) error
}
