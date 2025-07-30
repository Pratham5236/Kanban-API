
package middlewares

import (
	"net/http"

	"kanban-app/api/auth"
	"kanban-app/api/models"

	"github.com/gin-gonic/gin"
)

func CasbinMiddleware(paramName, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "User not found in context"})
			c.Abort()
			return
		}

		obj := c.Param(paramName)

		can, err := auth.NewAuthorizationService().Enforce(userID.(string), obj, act)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error checking authorization"})
			c.Abort()
			return
		}

		if !can {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}
