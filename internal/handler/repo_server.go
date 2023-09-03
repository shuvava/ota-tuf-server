package handler

import (
	"net/http"
	"strconv"

	"github.com/shuvava/ota-tuf-server/pkg/api"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
	"github.com/shuvava/ota-tuf-server/pkg/services"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"
	"github.com/shuvava/go-ota-svc-common/apperrors"

	"github.com/labstack/echo/v4"
)

const (
	pathVersion = "version"
	//PathRepoServerRepoWithNameSpaceResolver is the path to create / modify a key repository with resolving repositoryID from namespace
	PathRepoServerRepoWithNameSpaceResolver = "/user_repo"

	//PathRepoServerRepoContentWithNameSpaceResolver is the path to get current version of TUF key repository signed content with resolving repositoryID from namespace
	PathRepoServerRepoContentWithNameSpaceResolver = PathRepoServerRepoWithNameSpaceResolver + "/root.json"

	// PathRepoServerRepos is repo-server path returns list of repositories
	PathRepoServerRepos = "/repos"
)

type (
	rootGenRequest struct {
		Threshold uint               `json:"threshold,omitempty"`
		KeyType   encryption.KeyType `json:"keyType,omitempty"`
	}
)

func (r *rootGenRequest) Validate() error {
	if r.Threshold < 1 {
		return apperrors.NewAppError(apperrors.ErrorDataValidation, "incorrect threshold value: ")
	}
	return r.KeyType.Validate()
}

// ListRepos returns list of available repositories for all Namespaces
func ListRepos(ctx echo.Context, svc *services.RepositoryService) error { //nolint:typecheck
	c := cmnapi.GetRequestContext(ctx)
	offset, err := getInt64Param(ctx, "offset", 0)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	limit, err := getInt64Param(ctx, "limit", 10)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	res, total, err := svc.List(c, offset, limit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	m := api.PaginatedResponse[data.Repo]{
		Data:  res,
		Total: total,
	}
	return ctx.JSON(http.StatusOK, m)
}

func getInt64Param(ctx echo.Context, id string, defaultValue int64) (int64, error) {
	param := ctx.QueryParam(id)
	if param == "" {
		return defaultValue, nil
	}
	return strconv.ParseInt(param, 10, 64)
}
