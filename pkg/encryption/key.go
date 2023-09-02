package encryption

// BaseKey is interface with base ky functionality
type BaseKey interface {
	// Type returns the type of key.
	Type() KeyType
	// Method returns the method of signature
	Method() KeyMethod
	// Public this is the public string used as a unique identifier for the verifier instance.
	Public() string
	// FingerprintSHA256 returns the SHA256 fingerprint of the given key.
	FingerprintSHA256() string
}

// Key represents a common methods of different keys.
type Key interface {
	// MarshalAllData returns the data.SerializedKey object associated with the verifier contains public and private keys.
	MarshalAllData() (*SerializedKey, error)
	// MarshalPublicData returns the data.SerializedKey object associated with the verifier contains only public key.
	MarshalPublicData() (*SerializedKey, error)
	Signer
	Verifier
}
