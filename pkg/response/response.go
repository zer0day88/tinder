package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type ApiJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"user,omitempty"`
	Err     error  `json:"-"`
}

var (
	ErrSystem = ApiJSON{
		Message: "Unexpected error occurred, please try again later",
		Code:    http.StatusInternalServerError,
	}

	ErrNotFound = ApiJSON{
		Message: "Data Not Found",
		Code:    http.StatusNotFound,
	}

	ErrUnauthorized = ApiJSON{
		Message: "Request unauthorized. Missing or invalid access token",
		Code:    http.StatusUnauthorized,
	}

	ErrForbidden = ApiJSON{
		Message: "Current user can't request to change this current resource",
		Code:    http.StatusForbidden,
	}

	ErrBadRequest = ApiJSON{
		Message: "Invalid request params, header, or body",
		Code:    http.StatusBadRequest,
	}

	OKNoError = ApiJSON{
		Message: "Success",
		Code:    http.StatusOK,
	}
)

func (a ApiJSON) New(code int, msg string, data any) ApiJSON {
	return ApiJSON{
		Code:    code,
		Message: msg,
		Data:    data,
	}
}

func (a ApiJSON) WithData(data any) ApiJSON {
	a.Data = data
	return a
}

func (a ApiJSON) WithMsg(msg string) ApiJSON {
	a.Message = msg
	return a
}

func (a ApiJSON) WithErr(err error) ApiJSON {
	a.Err = err
	return a
}

func WriteJSON(c echo.Context, json ApiJSON) error {
	return c.JSON(json.Code, json)
}
