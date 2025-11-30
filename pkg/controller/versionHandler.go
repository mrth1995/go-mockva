package controller

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/mrth1995/go-mockva/pkg/version"
)

type VersionController struct{}

func NewVersionController() *VersionController {
	return &VersionController{}
}

type VersionResponse struct {
	Version string `json:"version"`
}

func (v *VersionController) GetVersion(request *restful.Request, response *restful.Response) {
	versionResp := VersionResponse{
		Version: version.Version,
	}
	response.WriteHeaderAndEntity(http.StatusOK, versionResp)
}

func (v *VersionController) RegisterEndpoint(ws *restful.WebService) {
	ws.Route(ws.GET("/version").
		To(v.GetVersion).
		Doc("Get application version").
		Returns(200, "OK", VersionResponse{}))
}
