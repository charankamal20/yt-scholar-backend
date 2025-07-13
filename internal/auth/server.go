// @title           Auth Service API
// @version         1.0
// @description     Handles Auth with Google OAuth authentication.
// @host            localhost:8080
// @BasePath        /auth

package auth

import (
	"os"

	authStore "github.com/charankamal20/youtube-scholar-backend/database/repository/auth"
	_ "github.com/charankamal20/youtube-scholar-backend/docs"
	common "github.com/charankamal20/youtube-scholar-backend/internal/common"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthServer struct {
	config  *oauth2.Config
	service AuthService
	*common.Server
}

func NewAuth(srv *common.Server, store authStore.Querier) *AuthServer {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SEC")
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/api/v1/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	authServer := &AuthServer{
		config:  conf,
		service: newAuthService(store),
		Server:  srv,
	}
	authServer.registerRoutes()
	return authServer
}

func (a *AuthServer) registerRoutes() {
	authServer := a.Public.Group("/auth")
	authServer.GET("/google/login", a.loginHandler)
	authServer.GET("/google/register", a.registerHandler)
	authServer.GET("/google/callback", a.oAuthCallbackHandler)
	authServer.GET("/public-key", a.getPublicKeyHandler)

	protectedAuthServer := a.Private.Group("/auth")
	protectedAuthServer.GET("/user", a.getUserInfoHandler)
}
