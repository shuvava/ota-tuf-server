package data

type (
	// KeyGenRequestStatus is a type of TUF server key generation
	KeyGenRequestStatus uint32
)

const (
	// KeyGenRequestStatusPending is pending status of key generation
	KeyGenRequestStatusPending KeyGenRequestStatus = iota
	// KeyGenRequestStatusSuccess is success status of key generation
	KeyGenRequestStatusSuccess
	// KeyGenRequestStatusFailed is an error status of key generation
	KeyGenRequestStatusFailed
)
