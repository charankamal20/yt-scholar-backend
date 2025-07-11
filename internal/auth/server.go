// @title           Auth Service API
// @version         1.0
// @description     Handles Auth with Google OAuth authentication.
// @host            localhost:8080
// @BasePath        /auth

package auth

import (
	"fmt"
	"os"

	authStore "github.com/charankamal20/youtube-scholar-backend/database/repository/auth"
	_ "github.com/charankamal20/youtube-scholar-backend/docs"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthServer struct {
	config *oauth2.Config
	server *gin.Engine

	service AuthService
}

func NewAuth(srv *gin.Engine, store authStore.Querier) *AuthServer {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SEC")
	fmt.Println(clientID, clientSecret)

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	authServer := &AuthServer{
		config:  conf,
		server:  srv,
		service: newAuthService(store),
	}

	authServer.registerRoutes()
	return authServer
}

func (a *AuthServer) registerRoutes() {
	authServer := a.server.Group("/auth")

	authServer.GET("/google/login", a.loginHandler)
	authServer.GET("/google/register", a.registerHandler)

	authServer.GET("/google/callback", a.oAuthCallbackHandler)
}
