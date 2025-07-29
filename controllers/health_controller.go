package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
// @Summary Show the status of the server.
// @Description get the status of the server.
// @Tags health
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{} "OK"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Kanban API is running!",
	})
}
