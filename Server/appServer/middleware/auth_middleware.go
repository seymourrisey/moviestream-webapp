package middleware

import (
	"net/http"

	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Token Provided"})
			c.Abort()
			return
		}
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}
		c.Set("userId", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
