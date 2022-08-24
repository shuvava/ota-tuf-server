package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"
	"github.com/shuvava/go-ota-svc-common/apperrors"

	"github.com/shuvava/ota-tuf-server/internal/db/mongo"
	tufapi "github.com/shuvava/ota-tuf-server/pkg/api"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
	"github.com/shuvava/ota-tuf-server/pkg/services"
)

const (
	pathRepoID = "repoID"
	//PathRepoServerRepoWithNameSpaceResolver is the path to create / modify a key repository with resolving repositoryID from namespace
	PathRepoServerRepoWithNameSpaceResolver = "/user_repo"
	//PathRepoServerRepo is the path to create/ modify a key repository
	PathRepoServerRepo = "/repo/:" + pathRepoID
	//PathKeyServerRepo is the path to create a new key repository
	PathKeyServerRepo = "/root/:" + pathRepoID
	// PathRepoServerRepos is repo-server path returns list of repositories
	PathRepoServerRepos = "/repos"
)

type (
	rootGenRequest struct {
		Threshold int                `json:"threshold,omitempty"`
		KeyType   encryption.KeyType `json:"keyType,omitempty"`
	}
)

// CreateRoot creates a new TUF key repository
func CreateRoot(ctx echo.Context, svc *services.RepositoryService) error {
	c := cmnapi.GetRequestContext(ctx)
	ns := cmnapi.GetNamespace(ctx)
	repoID, err := getRepoID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	genReq := &rootGenRequest{
		Threshold: 1,
		KeyType:   encryption.KeyTypeRSA,
	}
	if err = ctx.Bind(genReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	err = svc.Create(c, ns, repoID, genReq.KeyType)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == mongo.ErrorRepoErrorDbAlreadyExist {
			return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
		}
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(tufapi.HeaderRepoID, repoID.String())
	return ctx.NoContent(http.StatusOK)
}

// ListRepos returns list of available repositories for all Namespaces
func ListRepos(ctx echo.Context, svc *services.RepositoryService) error {
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
	m := tufapi.PaginatedResponse[data.Repo]{
		Data:  res,
		Total: total,
	}
	return ctx.JSON(http.StatusOK, m)
}

func getRepoID(ctx echo.Context) (data.RepoID, error) {
	repoID := ctx.Param(pathRepoID)
	if repoID == "" {
		if strings.HasSuffix(ctx.Path(), PathRepoServerRepoWithNameSpaceResolver) {
			// TODO: make repoID generation consistent (UUIDv5 (namespace_name)
			return data.NewRepoID(), nil
		}
		return data.RepoIDNil, apperrors.NewAppError(apperrors.ErrorGeneric, "parameter repoID is missing")
	}
	return data.RepoIDFromString(repoID)
}

func getInt64Param(ctx echo.Context, id string, defaultValue int64) (int64, error) {
	param := ctx.QueryParam(id)
	if param == "" {
		return defaultValue, nil
	}
	return strconv.ParseInt(param, 10, 64)
}
