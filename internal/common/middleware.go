package common

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/charankamal20/youtube-scholar-backend/pkg/token"
	"github.com/gin-gonic/gin"
)

const (
	// Context keys for storing token data
	AuthUserIDKey  = "auth_user_id"
	AuthEmailKey   = "auth_email"
	AuthPayloadKey = "auth_payload"
	AuthRolesKey   = "auth_roles"
	AuthSubjectKey = "auth_subject"
	AuthTokenKey   = "auth_token"

	// Header keys
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
)

// AuthMiddlewareConfig holds configuration for the auth middleware
type AuthMiddlewareConfig struct {
	TokenMaker    *token.PasetoMaker
	SkipPaths     []string
	RequiredRoles []string
}

func AuthMiddleware(config AuthMiddlewareConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if shouldSkipAuth(ctx.Request.URL.Path, config.SkipPaths) {
			ctx.Next()
			return
		}

		accessToken, err := ctx.Cookie("access_token")
		if err != nil || accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := config.TokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if len(config.RequiredRoles) > 0 {
			if !hasRequiredRole(payload.Roles, config.RequiredRoles) {
				err := errors.New("insufficient permissions")
				ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.Set(AuthUserIDKey, payload.UserID)
		ctx.Set(AuthEmailKey, payload.Email)
		ctx.Set(AuthPayloadKey, payload)
		ctx.Set(AuthRolesKey, payload.Roles)
		ctx.Set(AuthSubjectKey, payload.Subject)
		ctx.Set(AuthTokenKey, accessToken)

		ctx.Next()
	}
}

func RequireAuth(tokenMaker *token.PasetoMaker) gin.HandlerFunc {
	return AuthMiddleware(AuthMiddlewareConfig{
		TokenMaker: tokenMaker,
	})
}

func RequireRoles(tokenMaker *token.PasetoMaker, roles ...string) gin.HandlerFunc {
	return AuthMiddleware(AuthMiddlewareConfig{
		TokenMaker:    tokenMaker,
		RequiredRoles: roles,
	})
}

func shouldSkipAuth(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

func hasRequiredRole(userRoles, requiredRoles []string) bool {
	if len(requiredRoles) == 0 {
		return true
	}

	roleMap := make(map[string]bool)
	for _, role := range userRoles {
		roleMap[role] = true
	}

	for _, requiredRole := range requiredRoles {
		if roleMap[requiredRole] {
			return true
		}
	}
	return false
}

// errorResponse creates a standard error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func GetAuthUserID(ctx *gin.Context) (string, bool) {
	userID, exists := ctx.Get(AuthUserIDKey)
	if !exists {
		return "", false
	}
	return userID.(string), true
}

func GetAuthEmail(ctx *gin.Context) (string, bool) {
	email, exists := ctx.Get(AuthEmailKey)
	if !exists {
		return "", false
	}
	return email.(string), true
}

func GetAuthPayload(ctx *gin.Context) (*token.Payload, bool) {
	payload, exists := ctx.Get(AuthPayloadKey)
	if !exists {
		return nil, false
	}
	return payload.(*token.Payload), true
}

func GetAuthRoles(ctx *gin.Context) ([]string, bool) {
	roles, exists := ctx.Get(AuthRolesKey)
	if !exists {
		return nil, false
	}
	return roles.([]string), true
}

func GetAuthSubject(ctx *gin.Context) (string, bool) {
	subject, exists := ctx.Get(AuthSubjectKey)
	if !exists {
		return "", false
	}
	return subject.(string), true
}

func GetAuthToken(ctx *gin.Context) (string, bool) {
	token, exists := ctx.Get(AuthTokenKey)
	if !exists {
		return "", false
	}
	return token.(string), true
}

func HasRole(ctx *gin.Context, role string) bool {
	roles, exists := GetAuthRoles(ctx)
	if !exists {
		return false
	}

	for _, userRole := range roles {
		if userRole == role {
			return true
		}
	}
	return false
}

func HasAnyRole(ctx *gin.Context, roles ...string) bool {
	userRoles, exists := GetAuthRoles(ctx)
	if !exists {
		return false
	}

	return hasRequiredRole(userRoles, roles)
}

func IsAuthenticated(ctx *gin.Context) bool {
	_, exists := GetAuthUserID(ctx)
	return exists
}

// Advanced middleware variants

// RateLimitByUser creates a rate limiter based on user ID
func RateLimitByUser(tokenMaker *token.PasetoMaker, maxRequests int, window time.Duration) gin.HandlerFunc {
	// This would typically use Redis or similar for distributed rate limiting
	// For simplicity, using in-memory map (not suitable for production)
	userLimits := make(map[string]time.Time)
	mu := sync.RWMutex{}

	return gin.HandlerFunc(func(ctx *gin.Context) {
		// First verify the token
		authMiddleware := RequireAuth(tokenMaker)
		authMiddleware(ctx)

		if ctx.IsAborted() {
			return
		}

		userID, exists := GetAuthUserID(ctx)
		if !exists {
			ctx.Next()
			return
		}

		mu.Lock()
		defer mu.Unlock()

		// Simple rate limiting logic (extend as needed)
		if lastRequest, exists := userLimits[userID]; exists {
			if time.Since(lastRequest) < window {
				ctx.AbortWithStatusJSON(http.StatusTooManyRequests,
					gin.H{"error": "rate limit exceeded"})
				return
			}
		}

		userLimits[userID] = time.Now()
		ctx.Next()
	})
}
