package data

// SignedPayload is public model with signed content
type SignedPayload[T RootRole] struct {
	Signature *ClientSignature
	Signed    *T
}
