package handler

import (
	"strconv"
	"strings"

	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/errcodes"

	"github.com/shuvava/go-ota-svc-common/apperrors"

	"github.com/labstack/echo/v4"
)

func getOrGenerateRepoID(ctx echo.Context) (data.RepoID, error) {
	repoID, err := getRepoID(ctx)
	if err != nil {
		if strings.HasSuffix(ctx.Path(), PathRepoServerRepoWithNameSpaceResolver) {
			// TODO: make repoID generation consistent (UUIDv5 (namespace_name)
			return data.NewRepoID(), nil
		}
		return data.RepoIDNil, err
	}

	return repoID, nil
}

func getRepoID(ctx echo.Context) (data.RepoID, error) {
	k := ctx.Param(pathRepoID)
	if k == "" {
		return data.RepoIDNil, apperrors.NewAppError(errcodes.ErrorAPIRequestValidationParamMissing, "parameter repoID is missing")
	}
	repoID, err := data.RepoIDFromString(k)
	if err != nil {
		return data.RepoIDNil, apperrors.NewAppError(errcodes.ErrorAPIRequestValidationParamInvalid, "parameter repoID is invalid")
	}
	return repoID, nil
}

func getVersion(ctx echo.Context) (uint, error) {
	ver := ctx.Param(pathVersion)
	if ver == "" {
		return 0, apperrors.NewAppError(errcodes.ErrorAPIRequestValidationParamMissing, "parameter version is missing")
	}
	v, err := strconv.ParseUint(ver, 10, 32)
	if err != nil {
		return 0, apperrors.NewAppError(errcodes.ErrorAPIRequestValidationParamInvalid, "parameter version is invalid")
	}
	return uint(v), nil
}

func getKeyID(ctx echo.Context) (data.KeyID, error) {
	k := ctx.Param(pathKeyID)
	if k == "" {
		return "", apperrors.NewAppError(errcodes.ErrorAPIRequestValidationParamMissing, "parameter version is missing")
	}
	return data.KeyID(k), nil
}

func getRole(ctx echo.Context) (data.RoleType, error) {
	k := ctx.Param(pathRole)
	if k == "" {
		return "", apperrors.NewAppError(errcodes.ErrorAPIRequestValidationParamMissing, "parameter role is missing")
	}
	role, err := data.NewRoleType(k)
	if err != nil {
		return "", apperrors.NewAppError(errcodes.ErrorAPIRequestValidationParamInvalid, "parameter role is invalid")
	}
	return role, nil
}
