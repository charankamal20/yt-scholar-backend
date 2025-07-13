// Updated common/server.go
package common

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/charankamal20/youtube-scholar-backend/pkg/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine       *gin.Engine
	Public       *gin.RouterGroup
	Private      *gin.RouterGroup   // Changed from gin.IRoutes to *gin.RouterGroup
	tokenService *token.PasetoMaker // Private - hidden from other packages
}

func NewServer() *Server {
	tokenService, err := token.NewPasetoMaker()
	if err != nil {
		log.Fatal("could not make token service")
		return nil
	}

	server := gin.Default()
	gin.SetMode(gin.DebugMode)

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
	}))

	serverGroup := server.Group("/api/v1")
	serverGroup.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})
	serverGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create a separate group for private routes with auth middleware
	privateGroup := serverGroup.Group("")
	privateGroup.Use(RequireAuth(tokenService))

	appServer := &Server{
		engine:       server,
		Public:       serverGroup,
		Private:      privateGroup,
		tokenService: tokenService,
	}

	return appServer
}

// Expose only the methods you need - token service remains private
func (s *Server) CreateToken(userID, email string, opts *token.TokenOptions) (string, error) {
	return s.tokenService.CreateToken(userID, email, opts)
}

func (s *Server) GetPublicKey() []byte {
	return s.tokenService.PublicKey()
}

// Optional: If you need to apply auth middleware manually somewhere
func (s *Server) GetAuthMiddleware() gin.HandlerFunc {
	return RequireAuth(s.tokenService)
}

func (s *Server) Run() {
	if err := s.engine.Run(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	fmt.Println("Started server on port 8080.")
}
