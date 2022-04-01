package middleware

import (
	"net/http"
	"os"
	"quik/internal/encryption"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthPlayer() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.GetHeader("Authorization")
		if clientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Authorization header provided"})
			c.Abort()
			return
		}

		idTokenHeader := strings.Split(clientToken, "Bearer ")

		if len(idTokenHeader) < 2 {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Must provide Authorization header with format `Bearer {token}`",
			})
			c.Abort()
			return
		}
		claims, err := encryption.ValidateToken(idTokenHeader[1], os.Getenv("SECRET_KEY"))
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}
		c.Set("playerId", claims.Id)
		c.Next()
	}
}
