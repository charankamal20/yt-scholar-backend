// @title           Auth Service API
// @version         1.0
// @description     Handles Auth with Google OAuth authentication.
// @host            localhost:8080
// @BasePath        /auth

package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	authStore "github.com/charankamal20/youtube-scholar-backend/database/repository/auth"
	_ "github.com/charankamal20/youtube-scholar-backend/docs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthServer struct {
	config *oauth2.Config
	server *gin.Engine

	store *authStore.Querier
}

func NewAuth(srv *gin.Engine, store authStore.Querier) *AuthServer {
	err := godotenv.Load()

	if err != nil {
		log.Println("No env file found!")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SEC")

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	authServer := &AuthServer{
		config: conf,
		server: srv,
		store:  &store,
	}

	authServer.registerRoutes()
	return authServer
}

func (a *AuthServer) registerRoutes() {
	authServer := a.server.Group("/auth")

	authServer.GET("/google/login", a.oAuthHandler)
	authServer.GET("/google/callback", a.oAuthCallbackHandler)
}

// oAuthHandler redirects the user to Google's OAuth page
func (a *AuthServer) oAuthHandler(c *gin.Context) {
	url := a.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// oAuthCallbackHandler handles the Google callback
// @Summary      Google OAuth Callback
// @Description  Handles the Google OAuth 2.0 callback, exchanges code for token, and fetches user info.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        code  query     string  true  "OAuth2 Authorization Code"
// @Success      200   {object}  map[string]interface{}  "User info from Google"
// @Failure      400   {object}  map[string]interface{}  "Bad request or token exchange failure"
// @Failure      500   {object}  map[string]interface{}  "Internal server error"
// @Router       /auth/google/callback [get]
func (a *AuthServer) oAuthCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
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

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
