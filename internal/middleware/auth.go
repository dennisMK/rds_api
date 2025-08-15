package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"healthcare-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	jwtSecret []byte
	logger    *logrus.Logger
}

func NewAuthMiddleware(jwtSecret string, logger *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(jwtSecret),
		logger:    logger,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Scopes   []string `json:"scopes"`
	jwt.RegisteredClaims
}

// RequireAuth middleware validates JWT tokens
func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			a.logger.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, models.NewOperationOutcome("error", "security", "Authorization header required"))
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			a.logger.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, models.NewOperationOutcome("error", "security", "Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims := &Claims{}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return a.jwtSecret, nil
		})

		if err != nil {
			a.logger.WithError(err).Warn("Invalid JWT token")
			c.JSON(http.StatusUnauthorized, models.NewOperationOutcome("error", "security", "Invalid or expired token"))
			c.Abort()
			return
		}

		if !token.Valid {
			a.logger.Warn("Invalid JWT token")
			c.JSON(http.StatusUnauthorized, models.NewOperationOutcome("error", "security", "Invalid token"))
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("roles", claims.Roles)
		c.Set("scopes", claims.Scopes)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func (a *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			a.logger.Error("Roles not found in context")
			c.JSON(http.StatusForbidden, models.NewOperationOutcome("error", "security", "Access denied"))
			c.Abort()
			return
		}

		userRoles, ok := roles.([]string)
		if !ok {
			a.logger.Error("Invalid roles format in context")
			c.JSON(http.StatusForbidden, models.NewOperationOutcome("error", "security", "Access denied"))
			c.Abort()
			return
		}

		// Check if user has required role
		hasRole := false
		for _, role := range userRoles {
			if role == requiredRole || role == "admin" { // admin has access to everything
				hasRole = true
				break
			}
		}

		if !hasRole {
			a.logger.WithFields(logrus.Fields{
				"required_role": requiredRole,
				"user_roles":    userRoles,
			}).Warn("Insufficient permissions")
			c.JSON(http.StatusForbidden, models.NewOperationOutcome("error", "security", "Insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireScope middleware checks if user has required scope
func (a *AuthMiddleware) RequireScope(requiredScope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopes, exists := c.Get("scopes")
		if !exists {
			a.logger.Error("Scopes not found in context")
			c.JSON(http.StatusForbidden, models.NewOperationOutcome("error", "security", "Access denied"))
			c.Abort()
			return
		}

		userScopes, ok := scopes.([]string)
		if !ok {
			a.logger.Error("Invalid scopes format in context")
			c.JSON(http.StatusForbidden, models.NewOperationOutcome("error", "security", "Access denied"))
			c.Abort()
			return
		}

		// Check if user has required scope
		hasScope := false
		for _, scope := range userScopes {
			if scope == requiredScope || scope == "*" { // * grants all scopes
				hasScope = true
				break
			}
		}

		if !hasScope {
			a.logger.WithFields(logrus.Fields{
				"required_scope": requiredScope,
				"user_scopes":    userScopes,
			}).Warn("Insufficient scope")
			c.JSON(http.StatusForbidden, models.NewOperationOutcome("error", "security", "Insufficient scope"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateToken generates a JWT token for a user
func (a *AuthMiddleware) GenerateToken(userID, username string, roles, scopes []string, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		Scopes:   scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "healthcare-api",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

// GetUserFromContext extracts user information from gin context
func GetUserFromContext(c *gin.Context) (userID, username string, roles, scopes []string) {
	if uid, exists := c.Get("user_id"); exists {
		userID, _ = uid.(string)
	}
	if uname, exists := c.Get("username"); exists {
		username, _ = uname.(string)
	}
	if r, exists := c.Get("roles"); exists {
		roles, _ = r.([]string)
	}
	if s, exists := c.Get("scopes"); exists {
		scopes, _ = s.([]string)
	}
	return
}
