package handler

import (
	"strconv"
	"strings"

	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/errcodes"

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
	repoID := ctx.Param(pathRepoID)
	if repoID == "" {
		return data.RepoIDNil, apperrors.NewAppError(errcodes.ErrorAPIRequestValidation, "parameter repoID is missing")
	}
	return data.RepoIDFromString(repoID)
}

func getVersion(ctx echo.Context) (uint, error) {
	ver := ctx.Param(pathVersion)
	if ver == "" {
		return 0, apperrors.NewAppError(errcodes.ErrorAPIRequestValidation, "parameter version is missing")
	}
	v, err := strconv.ParseUint(ver, 10, 32)
	if err != nil {
		return 0, apperrors.NewAppError(errcodes.ErrorAPIRequestValidation, "parameter version is invalid")
	}
	return uint(v), nil
}
