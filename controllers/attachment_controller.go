
package controllers

import (
	"net/http"

	"kanban-app/api/models"
	"kanban-app/api/services"

	"github.com/gin-gonic/gin"
)

var attachmentService *services.AttachmentService

func init() {
	attachmentService = services.NewAttachmentService()
}

// CreateAttachment handles creating a new attachment on a card.
// @Summary Create a new attachment
// @Description Creates a new attachment on a specified card.
// @Tags Attachments
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param cardID path string true "Card ID"
// @Param attachment body models.CreateAttachmentRequest true "Attachment creation details"
// @Success 201 {object} models.Attachment "Attachment created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /cards/{cardID}/attachments [post]
func CreateAttachment(c *gin.Context) {
	cardID := c.Param("cardID")

	var req models.CreateAttachmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	attachment, err := attachmentService.CreateAttachment(cardID, req.FileName, req.FileURL, req.FileType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to create attachment: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

// DeleteAttachment handles deleting an attachment.
// @Summary Delete an attachment
// @Description Deletes a specific attachment by its ID.
// @Tags Attachments
// @Security ApiKeyAuth
// @Param cardID path string true "Card ID"
// @Param attachmentID path string true "Attachment ID"
// @Success 204 "No Content"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /cards/{cardID}/attachments/{attachmentID} [delete]
func DeleteAttachment(c *gin.Context) {
	attachmentID := c.Param("attachmentID")

	if err := attachmentService.DeleteAttachment(attachmentID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to delete attachment: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
