// @title           Auth Service API
// @version         1.0
// @description     Handles Auth with Google OAuth authentication.
// @host            localhost:8080
// @BasePath        /auth

package auth

import (
	"log"
	"os"

	authStore "github.com/charankamal20/youtube-scholar-backend/database/repository/auth"
	_ "github.com/charankamal20/youtube-scholar-backend/docs"
	middleware "github.com/charankamal20/youtube-scholar-backend/internal/common"
	"github.com/charankamal20/youtube-scholar-backend/pkg/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthServer struct {
	config       *oauth2.Config
	server       *gin.Engine
	tokenService *token.PasetoMaker

	service AuthService
}

func NewAuth(srv *gin.Engine, store authStore.Querier) *AuthServer {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SEC")

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	tokenService, err := token.NewPasetoMaker()
	if err != nil {
		log.Fatal("could not make token service")
		return nil
	}

	authServer := &AuthServer{
		config:       conf,
		server:       srv,
		service:      newAuthService(store),
		tokenService: tokenService,
	}

	authServer.registerRoutes()
	return authServer
}

func (a *AuthServer) registerRoutes() {
	authServer := a.server.Group("/auth")

	authMiddleware := middleware.RequireAuth(a.tokenService)

	authServer.GET("/google/login", a.loginHandler)
	authServer.GET("/google/register", a.registerHandler)

	authServer.GET("/google/callback", a.oAuthCallbackHandler)

	authServer.GET("/public-key", a.getPublicKeyHandler)

	protectedAuthServer := authServer.Use(authMiddleware)
	protectedAuthServer.GET("/user", authMiddleware, a.getUserInfoHandler)
}
