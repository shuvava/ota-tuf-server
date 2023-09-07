package data

// SignedPayload is public model with signed content
type SignedPayload[T RepoSigned | RoleSign] struct {
	Signatures []*ClientSignature `json:"signatures,omitempty"`
	Signed     *T                 `json:"signed"`
}

// RoleSign is a base type of allowed roles
type RoleSign struct {
	Role      RoleType `json:"keyType"`
	Threshold uint     `json:"threshold"`
}

// NewSignedPayload sings payload with RoleType
func NewSignedPayload(payload interface{}, role RoleType, keys []*RepoKey, threshold uint) (*SignedPayload[RoleSign], error) {
	signatures, err := SignPayload(payload, keys)
	if err != nil {
		return nil, err
	}

	sign := &RoleSign{
		Role:      role,
		Threshold: threshold,
	}
	return &SignedPayload[RoleSign]{
		Signatures: signatures,
		Signed:     sign,
	}, nil
}
