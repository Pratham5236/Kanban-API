package controllers

import (
	"net/http"
	"strings"

	"kanban-app/api/models"
	"kanban-app/api/services"

	"github.com/gin-gonic/gin"
)

var labelService *services.LabelService

func init() {
	labelService = services.NewLabelService()
}

// CreateLabel handles creating a new label.
// @Summary Create a new label
// @Description Creates a new reusable label.
// @Tags Labels
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param label body models.CreateLabelRequest true "Label creation details"
// @Success 201 {object} models.Label "Label created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 409 {object} models.ErrorResponse "Conflict"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /labels [post]
func CreateLabel(c *gin.Context) {
	var req models.CreateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	label, err := labelService.CreateLabel(req.Name, req.Color)
	if err != nil {
		if strings.Contains(err.Error(), "label name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to create label: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, label)
}

// GetAllLabels handles retrieving all labels.
// @Summary Get all labels
// @Description Retrieves all reusable labels.
// @Tags Labels
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Label "List of labels"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /labels [get]
func GetAllLabels(c *gin.Context) {
	labels, err := labelService.GetAllLabels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve labels: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, labels)
}

// GetLabelByID handles retrieving a label by ID.
// @Summary Get label by ID
// @Description Retrieves a specific label by its ID.
// @Tags Labels
// @Security ApiKeyAuth
// @Param labelID path string true "Label ID"
// @Produce json
// @Success 200 {object} models.Label "Label details"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /labels/{labelID} [get]
func GetLabelByID(c *gin.Context) {
	labelID := c.Param("labelID")

	label, err := labelService.GetLabelByID(labelID)
	if err != nil {
		if strings.Contains(err.Error(), "label not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve label: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, label)
}

// UpdateLabel handles updating an existing label.
// @Summary Update a label
// @Description Updates a specific label by its ID.
// @Tags Labels
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param labelID path string true "Label ID"
// @Param label body models.UpdateLabelRequest true "Label update details"
// @Success 200 {object} models.Label "Label updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 409 {object} models.ErrorResponse "Conflict"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /labels/{labelID} [put]
func UpdateLabel(c *gin.Context) {
	labelID := c.Param("labelID")

	var req models.UpdateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	updatedLabel, err := labelService.UpdateLabel(labelID, req)
	if err != nil {
		if strings.Contains(err.Error(), "label name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to update label: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedLabel)
}

// DeleteLabel handles deleting a label.
// @Summary Delete a label
// @Description Deletes a specific label by its ID.
// @Tags Labels
// @Security ApiKeyAuth
// @Param labelID path string true "Label ID"
// @Success 204 "No Content"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /labels/{labelID} [delete]
func DeleteLabel(c *gin.Context) {
	labelID := c.Param("labelID")

	err := labelService.DeleteLabel(labelID)
	if err != nil {
		if strings.Contains(err.Error(), "label not found or already deleted") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to delete label: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
