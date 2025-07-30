package controllers

import (
	"net/http"
	"strings"

	"kanban-app/api/models"
	"kanban-app/api/services"

	"github.com/gin-gonic/gin"
)

var boardService *services.BoardService

func init() {
	boardService = services.NewBoardService()
}

// CreateBoard handles creating a new board within a project.
// @Summary Create a new board
// @Description Creates a new board within a specified project. User must own the project.
// @Tags Boards
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param board body models.CreateBoardRequest true "Board creation details"
// @Success 201 {object} models.Board "Board created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 409 {object} models.ErrorResponse "Conflict"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards [post]
func CreateBoard(c *gin.Context) {
	userID, _ := c.Get("userID")
	projectID := c.Param("projectID")

	var req models.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	board, err := boardService.CreateBoard(projectID, req.Name, req.Description, userID.(string))
	if err != nil {
		if strings.Contains(err.Error(), "board name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to create board: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, board)
}

// GetBoards handles retrieving all boards for a specific project.
// @Summary Get all boards in a project
// @Description Retrieves all boards within a specified project. User must own the project.
// @Tags Boards
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Produce json
// @Success 200 {array} models.Board "List of boards"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards [get]
func GetBoards(c *gin.Context) {
	projectID := c.Param("projectID")

	boards, err := boardService.GetBoardsByProjectID(projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve boards: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, boards)
}

// GetBoardByID handles retrieving a specific board within a project.
// @Summary Get board by ID
// @Description Retrieves a specific board by its ID within a specified project. User must own the project.
// @Tags Boards
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Produce json
// @Success 200 {object} models.Board "Board details"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID} [get]
func GetBoardByID(c *gin.Context) {
	boardID := c.Param("boardID")

	board, err := boardService.GetBoardByID(boardID)
	if err != nil {
		if strings.Contains(err.Error(), "board not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve board: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}

// UpdateBoard handles updating an existing board within a project.
// @Summary Update a board
// @Description Updates a specific board by its ID within a specified project. User must own the project.
// @Tags Boards
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param board body models.UpdateBoardRequest true "Board update details"
// @Success 200 {object} models.Board "Board updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 409 {object} models.ErrorResponse "Conflict"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID} [put]
func UpdateBoard(c *gin.Context) {
	boardID := c.Param("boardID")

	var req models.UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	updatedBoard, err := boardService.UpdateBoard(boardID, req)
	if err != nil {
		if strings.Contains(err.Error(), "board name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to update board: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedBoard)
}

// DeleteBoard handles deleting a board within a project.
// @Summary Delete a board
// @Description Deletes a specific board by its ID within a specified project. User must own the project.
// @Tags Boards
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Success 204 "No Content"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID} [delete]
func DeleteBoard(c *gin.Context) {
	boardID := c.Param("boardID")

	err := boardService.DeleteBoard(boardID)
	if err != nil {
		if strings.Contains(err.Error(), "board not found or already deleted") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to delete board: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetBoardDetails handles retrieving a full board with all its lists and cards.
// @Summary Get full board details
// @Description Retrieves a board and all of its nested lists, cards, labels, etc.
// @Tags Boards
// @Security ApiKeyAuth
// @Produce json
// @Param boardID path string true "Board ID"
// @Success 200 {object} models.Board "Full board details"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /boards/{boardID}/details [get]
func GetBoardDetails(c *gin.Context) {
	boardID := c.Param("boardID")

	board, err := boardService.GetBoardDetails(boardID)
	if err != nil {
		if strings.Contains(err.Error(), "board not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve board details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}
