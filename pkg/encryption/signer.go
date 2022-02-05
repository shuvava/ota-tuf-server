package encryption

// Signer is an interface for an opaque private key that can be used for signing operations.
type Signer interface {
	Key
	// SignMessage signs a message with the private key.
	SignMessage(message []byte) ([]byte, error)
}
