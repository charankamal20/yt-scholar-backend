package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	common "github.com/charankamal20/youtube-scholar-backend/internal/common"
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
// @Failure      400   {object}  HTTPError  "Bad request or token exchange failure"
// @Failure      500   {object}  HTTPError  "Internal server error"
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
		c.JSON(common.ErrBadRequest.Code, common.ErrBadRequest.GetGinError())
		return
	}

	token, err := a.CreateToken(userInfo.ID, userInfo.Email, nil)

	c.SetCookie("access_token", token, int(time.Hour*24), "/", "localhost", true, true)
	c.SetCookie("user_id", userInfo.ID, int(time.Hour*24), "/", "localhost", true, false)

	// c.JSON(http.StatusOK, "Logged In")
	c.Redirect(http.StatusPermanentRedirect, "http://localhost:3000/")
}

func (a *AuthServer) getUserInfoHandler(ctx *gin.Context) {
	userId, exists := common.GetAuthUserID(ctx)
	if !exists {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("user id not available"))
		return
	}

	user, err := a.service.getUserInfo(ctx.Request.Context(), userId)
	if err != nil {
		ctx.JSON(err.Code, err.GetGinError())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (a *AuthServer) getPublicKeyHandler(c *gin.Context) {
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Cache-Control", "public, max-age=3600")

	_, err := c.Writer.Write(a.GetPublicKey())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to serve public key",
		})
		return
	}
}
