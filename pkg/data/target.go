package data

import (
	"fmt"
	"strings"

	"github.com/shuvava/go-ota-svc-common/apperrors"
)

type (
	// TargetFormat is the format of TUF repository
	TargetFormat int
	// TargetName is the name of TUF repository
	TargetName string
	// TargetVersion is the version of TUF repository
	TargetVersion string
	// TargetFilename is the filename(item) in TUF repository
	TargetFilename string
)

// ErrorTargetFilenameValidation represents an error when the hardware identifier is invalid.
const ErrorTargetFilenameValidation = apperrors.ErrorDataValidation + ":TargetFilename"

const (
	// OSTREE target format of TUF repository in OSTRee format.
	OSTREE TargetFormat = iota
	// BINARY target format of TUF repository in binary format.
	BINARY
)

// Validate validates TargetFilename
func (tfn TargetFilename) Validate() error {
	tlen := len(tfn)
	if tlen == 0 || tlen > 254 || strings.Contains(string(tfn), "..") {
		return apperrors.NewAppError(
			ErrorTargetFilenameValidation,
			fmt.Sprintf("cannot be empty or bigger than 254 chars or contain `..` : %s", tfn))
	}
	return nil
}
