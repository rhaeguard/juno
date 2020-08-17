package auth

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// GetAppID extracts the AppID from the gin context
func GetAppID(c *gin.Context) string {
	claims := jwt.ExtractClaims(c)
	return claims["app_id"].(string)
}
