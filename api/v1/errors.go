package v1

import (
	"net/http"

	"github.com/danielkrainas/tinkersnest/api/errcode"
)

const ErrorGroup = "tinkersnest.api.v1"

var (
	ErrorCodeHookUnknown = errcode.Register(ErrorGroup, errcode.ErrorDescriptor{
		Value:          "HOOK_UNKNOWN",
		Message:        "hook not known to server",
		Description:    "This is returned if the hook ID used during an operation is unknown to the server.",
		HTTPStatusCode: http.StatusNotFound,
	})
)
