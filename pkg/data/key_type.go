package data

// KeyType is a string denoting a public key signature system
type KeyType string

const (
	// KeyTypeEd25519 is the type of Ed25519 keys.
	KeyTypeEd25519 = KeyType("ed25519")
	// KeyTypeECDSA is the type of ECDSA keys with SHA2 and P256.
	KeyTypeECDSA = KeyType("ecdsa")
	// KeyTypeRSA is the type of RSA keys with RSASSA-PSS and SHA256.
	KeyTypeRSA = KeyType("rsa")
)
