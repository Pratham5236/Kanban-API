package middlewares

import (
	"errors"
	"fmt"
	"kanban-app/api/models"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Invalid Authorization header format. Expected 'Bearer <token>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		jwtSecret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Token is expired"})
			} else if errors.Is(err, jwt.ErrSignatureInvalid) {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Invalid token signature"})
			} else if errors.Is(err, jwt.ErrTokenMalformed) {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Token is malformed"})
			} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Token is not valid yet"})
			} else {
				c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Invalid token: " + err.Error()})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "User ID not found in token claims"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

