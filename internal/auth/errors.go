package auth

import (
	"net/http"

	errors "github.com/charankamal20/youtube-scholar-backend/internal/common"
)

var (
	ErrUserNotFound = &errors.HTTPError{Code: http.StatusUnauthorized, Message: "User not found."}
	ErrEmailExists  = &errors.HTTPError{Code: http.StatusConflict, Message: "Email already registered."}
)
