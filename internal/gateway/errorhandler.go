package gateway

import (
	"net/http"
	"strings"
)

type Error struct {
	HTTPStatus int    `json:"http_status"`
	Message    string `json:"message"`
}

func newError(status int, statusMsg ...string) func(msg ...string) *Error {
	if len(statusMsg) == 0 {
		statusMsg = []string{http.StatusText(status)}
	}
	return func(msg ...string) *Error {
		return &Error{
			HTTPStatus: status,
			Message:    strings.Join(append(statusMsg, msg...), " "),
		}
	}
}

var (
	ErrInternalServerError = newError(http.StatusInternalServerError)
	ErrUnauthorized        = newError(http.StatusUnauthorized)
	ErrForbidden           = newError(http.StatusForbidden)
)

func ErrorResponse(err *Error, w http.ResponseWriter) {
	JSONResponse(w, err.HTTPStatus, err)
}
