package controller

import (
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	endpointError "github.com/mrth1995/go-mockva/pkg/errors"
	"github.com/mrth1995/go-mockva/pkg/model"
)

func (accountHandler *AccountController) RegisterEndpoint(ws *restful.WebService) {
	tags := []string{"Accounts"}
	ws.Route(
		ws.GET("/accounts/{accountId}").
			To(accountHandler.FindByUserID).
			Produces(restful.MIME_JSON).
			Param(restful.PathParameter("accountId", "Account ID")).
			Returns(http.StatusOK, "Account exist", model.AccountInfo{}).
			Returns(http.StatusBadRequest, "Validation error", endpointError.EndpointError{}).
			Returns(http.StatusNotFound, "Account not found", endpointError.EndpointError{}).
			Returns(http.StatusInternalServerError, "Internal server error", endpointError.EndpointError{}).
			Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	ws.Route(
		ws.POST("/accounts").
			To(accountHandler.CreateAccount).
			Consumes(restful.MIME_JSON).
			Produces(restful.MIME_JSON).
			Reads(model.AccountRegister{}).
			Returns(http.StatusOK, "Account successfully created", model.AccountInfo{}).
			Returns(http.StatusBadRequest, "Validation error", endpointError.EndpointError{}).
			Returns(http.StatusConflict, "Account already exist", endpointError.EndpointError{}).
			Returns(http.StatusInternalServerError, "Internal server error", endpointError.EndpointError{}).
			Metadata(restfulspec.KeyOpenAPITags, tags))

	ws.Route(
		ws.PATCH("/accounts/{accountId}").
			To(accountHandler.EditAccount).
			Consumes(restful.MIME_JSON).
			Produces(restful.MIME_JSON).
			Reads(model.AccountEdit{}).
			Returns(http.StatusOK, "Account successfully UPDATED", model.AccountInfo{}).
			Returns(http.StatusBadRequest, "Validation error", endpointError.EndpointError{}).
			Returns(http.StatusNotFound, "Account not exist", endpointError.EndpointError{}).
			Returns(http.StatusInternalServerError, "Internal server error", endpointError.EndpointError{}).
			Metadata(restfulspec.KeyOpenAPITags, tags))
}
