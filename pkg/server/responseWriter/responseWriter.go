package responseWriter

import (
	"net/http"

	endpointError "github.com/mrth1995/go-mockva/pkg/errors"

	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
)

func WriteOK(content interface{}, response *restful.Response) {
	err := response.WriteAsJson(content)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func WriteNoContent(response *restful.Response) {
	response.WriteHeader(http.StatusNoContent)
}

func WriteCreated(response *restful.Response) {
	response.WriteHeader(http.StatusCreated)
}

func WriteBadRequest(e error, response *restful.Response) {
	writeError(http.StatusBadRequest, e, response)
}

func WriteNotFound(e error, response *restful.Response) {
	writeError(http.StatusNotFound, e, response)
}

func WriteInternalServerError(e error, response *restful.Response) {
	writeError(http.StatusInternalServerError, e, response)
}

//TODO: wrap error into model
func writeError(httpStatus int, e error, response *restful.Response) {
	var resp *endpointError.EndpointError
	if endpointErr, ok := e.(*endpointError.EndpointError); ok {
		resp = endpointErr
	} else {
		resp = &endpointError.EndpointError{
			ErrorMessage: e.Error(),
			ErrorCode:    "96",
		}
	}

	err := response.WriteHeaderAndJson(httpStatus, resp, restful.MIME_JSON)
	if err != nil {
		logrus.Error(err)
		return
	}
}
