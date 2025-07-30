package controllers

import (
	"net/http"
	"strings"

	"kanban-app/api/models"
	"kanban-app/api/services"

	"github.com/gin-gonic/gin"
)

var listService *services.ListService

func init() {
	listService = services.NewListService()
}

// CreateList handles creating a new list within a board.
// @Summary Create a new list
// @Description Creates a new list within a specified board. User must own the board.
// @Tags Lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param list body models.CreateListRequest true "List creation details"
// @Success 201 {object} models.List "List created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 409 {object} models.ErrorResponse "Conflict"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists [post]
func CreateList(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardID := c.Param("boardID")

	var req models.CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	list, err := listService.CreateList(boardID, req.Name, userID.(string))
	if err != nil {
		if strings.Contains(err.Error(), "list name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to create list: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, list)
}

// GetLists handles retrieving all lists for a specific board.
// @Summary Get all lists in a board
// @Description Retrieves all lists within a specified board. User must own the board.
// @Tags Lists
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Produce json
// @Success 200 {array} models.List "List of lists"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists [get]
func GetLists(c *gin.Context) {
	boardID := c.Param("boardID")

	lists, err := listService.GetListsByBoardID(boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve lists: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, lists)
}

// GetListByID handles retrieving a specific list within a board.
// @Summary Get list by ID
// @Description Retrieves a specific list by its ID within a specified board. User must own the board.
// @Tags Lists
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param listID path string true "List ID"
// @Produce json
// @Success 200 {object} models.List "List details"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists/{listID} [get]
func GetListByID(c *gin.Context) {
	listID := c.Param("listID")

	list, err := listService.GetListByID(listID)
	if err != nil {
		if strings.Contains(err.Error(), "list not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve list: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, list)
}

// UpdateList handles updating an existing list within a board.
// @Summary Update a list
// @Description Updates a specific list by its ID within a specified board. User must own the board.
// @Tags Lists
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param listID path string true "List ID"
// @Param list body models.UpdateListRequest true "List update details"
// @Success 200 {object} models.List "List updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 409 {object} models.ErrorResponse "Conflict"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists/{listID} [put]
func UpdateList(c *gin.Context) {
	listID := c.Param("listID")

	var req models.UpdateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	if req.Position != nil {
		if err := listService.MoveList(listID, *req.Position); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to move list: " + err.Error()})
			return
		}
	}

	c.Status(http.StatusOK)
}

// DeleteList handles deleting a list within a board.
// @Summary Delete a list
// @Description Deletes a specific list by its ID within a specified board. User must own the board.
// @Tags Lists
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param listID path string true "List ID"
// @Success 204 "No Content"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists/{listID} [delete]
func DeleteList(c *gin.Context) {
	listID := c.Param("listID")

	err := listService.DeleteList(listID)
	if err != nil {
		if strings.Contains(err.Error(), "list not found or already deleted") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to delete list: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}