package handler

import (
	"errors"
	"net/http"

	"github.com/shuvava/ota-tuf-server/internal/db/mongo"
	"github.com/shuvava/ota-tuf-server/pkg/api"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
	"github.com/shuvava/ota-tuf-server/pkg/services"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"
	"github.com/shuvava/go-ota-svc-common/apperrors"

	"github.com/labstack/echo/v4"
)

// PathKeyServerRepo is the path to create a new key repository
const (
	PathKeyServerRepo = "/root/:" + pathRepoID
	//PathKeyServerRepoWithVersion is the path to create a new key repository
	PathKeyServerRepoWithVersion = PathKeyServerRepo + "/:" + pathVersion
)

// CreateRoot creates a new TUF key repository
func CreateRoot(ctx echo.Context, svc *services.RepositoryService) error { //nolint:typecheck
	c := cmnapi.GetRequestContext(ctx)
	ns := cmnapi.GetNamespace(ctx)
	repoID, err := getOrGenerateRepoID(ctx)
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
	if err := genReq.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	err = svc.Create(c, ns, repoID, genReq.KeyType)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && (typedErr.ErrorCode == mongo.ErrorRepoErrorDbAlreadyExist ||
			typedErr.ErrorCode == apperrors.ErrorDataValidation) {
			return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
		}
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(api.HeaderRepoID, repoID.String())
	return ctx.NoContent(http.StatusOK)
}
