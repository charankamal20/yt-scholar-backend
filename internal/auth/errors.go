package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPError struct {
	Code     int
	Message  string
	Internal error
}

func (e *HTTPError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Internal
}		

func (e *HTTPError) ResponseBody() map[string]string {
	return map[string]string{
		"error": e.Message,
	}
}

func (e *HTTPError) GetGinError() gin.H {
	log.Println(e.Internal)

	return gin.H{
		"error": e.Message,
	}
}

var (
	ErrUserNotFound = &HTTPError{Code: http.StatusUnauthorized, Message: "User not found."}
	ErrEmailExists  = &HTTPError{Code: http.StatusConflict, Message: "Email already registered."}

	ErrServer = func(err error) *HTTPError {
		return &HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal server error.",
			Internal: err,
		}
	}

	ErrBadRequest = &HTTPError{
		Code:    http.StatusBadRequest,
		Message: "Try Again!",
	}
)
