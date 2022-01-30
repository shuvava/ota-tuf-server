package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"

	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/services"
)

const (
	pathRepoID = "repoID"
	//PathCreateRoot is the path to create a new key repository
	PathCreateRoot = "/root/:" + pathRepoID
)

type (
	rootGenRequest struct {
		Threshold int          `json:"threshold,omitempty"`
		KeyType   data.KeyType `json:"keyType,omitempty"`
	}
)

// CreateRoot creates a new TUF key repository
func CreateRoot(ctx echo.Context, svc *services.RepositoryService) error {
	c := cmnapi.GetRequestContext(ctx)
	repoID, err := getRepoID(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	genReq := &rootGenRequest{
		Threshold: 1,
		KeyType:   data.KeyTypeRSA,
	}
	if err = ctx.Bind(genReq); err != nil {
		return ctx.JSON(http.StatusBadRequest, cmnapi.NewErrorResponse(c, http.StatusBadRequest, err))
	}
	err = svc.CreateNewRepository(c, repoID, genReq.KeyType)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, cmnapi.NewErrorResponse(c, http.StatusInternalServerError, err))
	}
	return ctx.NoContent(http.StatusOK)
}

func getRepoID(ctx echo.Context) (data.RepoID, error) {
	repoID := ctx.Param(pathRepoID)
	return data.RepoIDFromString(repoID)
}
