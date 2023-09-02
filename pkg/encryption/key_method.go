package encryption

// KeyMethod is a string denoting a public key signature system
type KeyMethod string

const (
	KeyMethodRsaPssSha256 = KeyMethod("rsassa-pss-sha256")
	KeyMethodED25519      = KeyMethod("ed25519")
	KeyMethodECPrime256V1 = KeyMethod("ecPrime256v1")
)
