package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

// loginHandler redirects the user to Google's OAuth page
func (a *AuthServer) loginHandler(c *gin.Context) {
	url := a.config.AuthCodeURL("login", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *AuthServer) registerHandler(c *gin.Context) {
	url := a.config.AuthCodeURL("register", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// oAuthCallbackHandler handles the Google callback
// @Summary      Google OAuth Callback
// @Description  Handles the Google OAuth 2.0 callback, exchanges code for token, and fetches user info.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        code  query     string  true  "OAuth2 Authorization Code"
// @Success      200   {object}  GoogleOAuthUser  "User info from Google"
// @Failure      400   {object}  HttpError  "Bad request or token exchange failure"
// @Failure      500   {object}  HttpError  "Internal server error"
// @Router       /auth/google/callback [get]
func (a *AuthServer) oAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})

		return
	}
	// Exchange code for token
	tok, err := a.config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token exchange failed", "details": err.Error()})
		return
	}

	// Create OAuth2 client with token
	client := a.config.Client(context.Background(), tok)

	// Request user info from Google
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user info", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	var userInfo GoogleOAuthUser
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	switch state {
	case "login":
		err := a.service.loginUser(c.Request.Context(), userInfo)
		if err != nil {
			c.JSON(err.Code, err.GetGinError())
			return
		}

	case "register":
		err := a.service.registerUser(c.Request.Context(), userInfo)
		if err != nil {
			c.JSON(err.Code, err.GetGinError())
			return
		}

	default:
		c.JSON(ErrBadRequest.Code, ErrBadRequest.GetGinError())
		return
	}

	c.JSON(http.StatusOK, "Logged In")
}
