package encryption

// KeyMethod is a string denoting a public key signature system
type KeyMethod string

const (
	// KeyMethodRsaPssSha256 default sign method for KeyTypeRSA
	KeyMethodRsaPssSha256 = KeyMethod("rsassa-pss-sha256")
	// KeyMethodED25519 default sign method for KeyTypeEd25519
	KeyMethodED25519 = KeyMethod("ed25519")
	// KeyMethodECPrime256V1 default sing method for KeyTypeECDSA
	KeyMethodECPrime256V1 = KeyMethod("ecPrime256v1")
)
