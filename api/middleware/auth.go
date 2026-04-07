package middleware

import (
	"net/http"
	"strings"

	"aubergine/internal/auth"
	"aubergine/internal/database"
	"aubergine/internal/models"

	"github.com/gin-gonic/gin"
)

// AuthRequired enforces JWT validation on protected routes
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Attach claims to request context
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

// MinimumTier enforces that a user has a specific subscription level or higher.
// For example, "premium" > "basic" > "free".
func MinimumTier(requiredTier string) gin.HandlerFunc {
	tierRanks := map[string]int{
		"free":    0,
		"basic":   1,
		"premium": 2,
	}

	return func(c *gin.Context) {
		userPlan, exists := c.Get("plan")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no valid active subscription tier found"})
			return
		}

		planStr, ok := userPlan.(string)
		if !ok || tierRanks[planStr] < tierRanks[requiredTier] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "upgrade your subscription tier to access this resource"})
			return
		}

		c.Next()
	}
}

// AdminRequired enforces that the user has admin role
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil || user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}

		c.Next()
	}
}
