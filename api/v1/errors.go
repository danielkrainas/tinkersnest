package v1

import (
	"net/http"

	"github.com/danielkrainas/gobag/api/errcode"
)

const ErrorGroup = "tinkersnest.api.v1"

var (
	ErrorCodeResourceUnknown = errcode.Register(ErrorGroup, errcode.ErrorDescriptor{
		Value:          "RESOURCE_UNKNOWN",
		Message:        "resource not known to server",
		Description:    "This is returned if the resource name used during an operation is unknown to the server.",
		HTTPStatusCode: http.StatusNotFound,
	})
)
