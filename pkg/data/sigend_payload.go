package data

// SignedPayload is public model with signed content
type SignedPayload[T RootRole] struct {
	Signatures []*ClientSignature `json:"signatures"`
	Signed     *T                 `json:"signed"`
}
