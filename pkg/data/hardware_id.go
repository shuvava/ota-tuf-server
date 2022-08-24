package data

import (
	"fmt"

	"github.com/shuvava/go-ota-svc-common/apperrors"
)

// HardwareIdentifier represents a hardware identifier.
type HardwareIdentifier string

// ErrorHardwareIdentifierValidation represents an error when the hardware identifier is invalid.
const ErrorHardwareIdentifierValidation = apperrors.ErrorDataValidation + ":HardwareIdentifier"

// Validate if HardwareIdentifier has valid format
func (h HardwareIdentifier) Validate() error {
	lenH := len(h)
	min := 0
	max := 200
	if lenH == min || lenH > max {
		return apperrors.NewAppError(
			ErrorHardwareIdentifierValidation,
			fmt.Sprintf("`%s` is not between %d and %d chars long", h, min, max))
	}
	return nil
}

// NewHardwareIdentifier creates a new HardwareIdentifier if str is valid
func NewHardwareIdentifier(str string) (HardwareIdentifier, error) {
	hid := HardwareIdentifier(str)
	if err := hid.Validate(); err != nil {
		return "", err
	}
	return hid, nil
}
