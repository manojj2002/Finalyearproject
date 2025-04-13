package utils

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // Ideally from env

func CreateJWT(username, userID string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"userId":   userID,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid auth header"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// You can optionally extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("username", claims["username"])
			c.Set("userId", claims["userId"])
		}

		c.Next()
	}
}

func GetUserIdFromContext(c *gin.Context) (string, bool) {
	userIdRaw, exists := c.Get("userId")
	if !exists {
		return "", false
	}
	userId, ok := userIdRaw.(string)
	if !ok {
		return "", false
	}
	return userId, true
}
