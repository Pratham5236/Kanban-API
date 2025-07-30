
package controllers

import (
	"net/http"

	"kanban-app/api/models"
	"kanban-app/api/services"

	"github.com/gin-gonic/gin"
)

var commentService *services.CommentService

func init() {
	commentService = services.NewCommentService()
}

// CreateComment handles creating a new comment on a card.
// @Summary Create a new comment
// @Description Creates a new comment on a specified card.
// @Tags Comments
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param cardID path string true "Card ID"
// @Param comment body models.CreateCommentRequest true "Comment creation details"
// @Success 201 {object} models.Comment "Comment created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /cards/{cardID}/comments [post]
func CreateComment(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardID := c.Param("cardID")

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	comment, err := commentService.CreateComment(cardID, userID.(string), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to create comment: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// DeleteComment handles deleting a comment.
// @Summary Delete a comment
// @Description Deletes a specific comment by its ID.
// @Tags Comments
// @Security ApiKeyAuth
// @Param cardID path string true "Card ID"
// @Param commentID path string true "Comment ID"
// @Success 204 "No Content"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /cards/{cardID}/comments/{commentID} [delete]
func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentID")

	if err := commentService.DeleteComment(commentID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to delete comment: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
