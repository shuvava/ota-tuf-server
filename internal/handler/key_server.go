package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shuvava/ota-tuf-server/internal/db/mongo"
	"github.com/shuvava/ota-tuf-server/pkg/api"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
	"github.com/shuvava/ota-tuf-server/pkg/errcodes"
	"github.com/shuvava/ota-tuf-server/pkg/services"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"
	"github.com/shuvava/go-ota-svc-common/apperrors"

	"github.com/labstack/echo/v4"
)

// PathKeyServerRepo is the path to create a new key repository
const (
	pathKeyID         = "keyID"
	pathRepoID        = "repoID"
	pathVersion       = "version"
	pathRole          = "role"
	PathKeyServerRepo = "/root/:" + pathRepoID
	//PathKeyServerRepoWithVersion is the path to create a new key repository
	PathKeyServerRepoWithVersion = PathKeyServerRepo + "/:" + pathVersion
	// PathKeyServerRepoWithKeyID is path to delete private key from the repo
	PathKeyServerRepoWithKeyID = PathKeyServerRepo + "/private_keys/:" + pathKeyID
	// PathKeyServerRepoWithRole is path to sign payload by provided role of the repo
	PathKeyServerRepoWithRole = PathKeyServerRepo + "/:" + pathRole
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
	genReq.KeyType = encryption.KeyTypeFromString(string(genReq.KeyType))
	if err = genReq.Validate(); err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	keys, err := svc.Create(c, ns, repoID, genReq.KeyType, genReq.Threshold)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && (typedErr.ErrorCode == mongo.ErrorRepoErrorDbAlreadyExist ||
			typedErr.ErrorCode == apperrors.ErrorDataValidation) {
			return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
		}
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(api.HeaderRepoID, repoID.String())
	return ctx.JSON(http.StatusOK, keys)
}

// GetRepoSignedContent returns current repo signed metadata
func GetRepoSignedContent(ctx echo.Context, svc *services.RepositoryService) error { //nolint:typecheck
	c := cmnapi.GetRequestContext(ctx)
	repoID, err := getRepoID(ctx)
	if repoID == data.RepoIDNil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == errcodes.ErrorAPIRequestValidationParamMissing {
			// try to get repoID from the namespace
			ns := cmnapi.GetNamespace(ctx)
			r, e := svc.FindByNamespace(c, ns)
			if e == nil {
				repoID = r.RepoID
				err = nil
			}
		}
	}
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	res, err := svc.GetAndRefresh(c, repoID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(api.HeaderRepoID, repoID.String())
	return ctx.JSON(http.StatusOK, res)
}

// GetRepoSignedContentForVersion returns repo signed metadata for requested version
func GetRepoSignedContentForVersion(ctx echo.Context, svc *services.RepositoryService, scsvc *services.RepoVersionService) error {
	c := cmnapi.GetRequestContext(ctx)
	ver, err := getVersion(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	repoID, err := getRepoID(ctx)
	if repoID == data.RepoIDNil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == errcodes.ErrorAPIRequestValidationParamMissing {
			// try to get repoID from the namespace
			ns := cmnapi.GetNamespace(ctx)
			r, e := svc.FindByNamespace(c, ns)
			if e == nil {
				repoID = r.RepoID
				err = nil
			}
		}
	}
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	res, err := scsvc.GetVersion(c, repoID, ver)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(api.HeaderRepoID, repoID.String())
	return ctx.JSON(http.StatusOK, res)
}

// DeletePrivateKey handler delete private key from TUF repo key
func DeletePrivateKey(ctx echo.Context, svc *services.KeyRepositoryService) error {
	c := cmnapi.GetRequestContext(ctx)
	repoID, err := getRepoID(ctx)
	if repoID == data.RepoIDNil || err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	keyID, err := getKeyID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	err = svc.DeletePrivateKey(c, repoID, keyID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	return ctx.NoContent(http.StatusNoContent)
}

// SignPayload handler to sign provided payload by role keys
func SignPayload(ctx echo.Context, svc *services.RepositoryService) error {
	c := cmnapi.GetRequestContext(ctx)
	repoID, err := getRepoID(ctx)
	if repoID == data.RepoIDNil || err != nil {
		err = apperrors.NewAppError(errcodes.ErrorAPIRequestValidationParamMissing, "parameter repoID missing or invalid")
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	role, err := getRole(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	var payload interface{}
	if err = json.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		err = apperrors.NewAppError(errcodes.ErrorAPIRequestValidationBodyValidation, "error on body validation")
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	res, err := svc.SignPayload(c, repoID, role, payload)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	ctx.Response().Header().Set(api.HeaderRepoID, repoID.String())
	return ctx.JSON(http.StatusOK, res)

}
