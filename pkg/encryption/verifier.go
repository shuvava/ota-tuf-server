package encryption

// A Verifier verifies public key signatures.
type Verifier interface {
	BaseKey
	// Verify takes a message and signature, all as byte slices,
	// and determines whether the signature is valid for the given
	// key and message.
	Verify(msg, sig []byte) error
}
