package controllers

import (
	"net/http"
	"strings"

	"kanban-app/api/models"
	"kanban-app/api/services"

	"github.com/gin-gonic/gin"
)

var cardService *services.CardService

func init() {
	cardService = services.NewCardService()
}

// CreateCard handles creating a new card within a list.
// @Summary Create a new card
// @Description Creates a new card within a specified list. User must own the list.
// @Tags Cards
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param listID path string true "List ID"
// @Param card body models.CreateCardRequest true "Card creation details"
// @Success 201 {object} models.Card "Card created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists/{listID}/cards [post]
func CreateCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	listID := c.Param("listID")

	var req models.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	card, err := cardService.CreateCard(
		listID,
		req.Title,
		req.Description,
		req.DueDate,
		userID.(string),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to create card: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, card)
}

// GetCards handles retrieving all cards for a specific list.
// @Summary Get all cards in a list
// @Description Retrieves all cards within a specified list. User must own the list.
// @Tags Cards
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param listID path string true "List ID"
// @Produce json
// @Success 200 {array} models.Card "List of cards"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists/{listID}/cards [get]
func GetCards(c *gin.Context) {
	listID := c.Param("listID")

	cards, err := cardService.GetCardsByListID(listID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve cards: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, cards)
}

// GetCardByID handles retrieving a specific card within a list.
// @Summary Get card by ID
// @Description Retrieves a specific card by its ID within a specified list. User must own the list.
// @Tags Cards
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param listID path string true "List ID"
// @Param cardID path string true "Card ID"
// @Produce json
// @Success 200 {object} models.Card "Card details"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists/{listID}/cards/{cardID} [get]
func GetCardByID(c *gin.Context) {
	cardID := c.Param("cardID")

	card, err := cardService.GetCardByID(cardID)
	if err != nil {
		if strings.Contains(err.Error(), "card not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve card: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}

// UpdateCard handles updating an existing card within a list.
// @Summary Update a card
// @Description Updates a specific card by its ID within a specified list. User must own the list. Supports moving card to another list.
// @Tags Cards
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param listID path string true "List ID"
// @Param cardID path string true "Card ID"
// @Param card body models.UpdateCardRequest true "Card update details"
// @Success 200 {object} models.Card "Card updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists/{listID}/cards/{cardID} [put]
func UpdateCard(c *gin.Context) {
	cardID := c.Param("cardID")

	var req models.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	if req.ListID != "" && req.Position != nil {
		if err := cardService.MoveCard(cardID, req.ListID, *req.Position); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to move card: " + err.Error()})
			return
		}
	}

	c.Status(http.StatusOK)
}

// DeleteCard handles deleting a card within a list.
// @Summary Delete a card
// @Description Deletes a specific card by its ID within a specified list. User must own the list.
// @Tags Cards
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param boardID path string true "Board ID"
// @Param listID path string true "List ID"
// @Param cardID path string true "Card ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID}/boards/{boardID}/lists/{listID}/cards/{cardID} [delete]
func DeleteCard(c *gin.Context) {
	cardID := c.Param("cardID")

	err := cardService.DeleteCard(cardID)
	if err != nil {
		if strings.Contains(err.Error(), "card not found or already deleted") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to delete card: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AddLabelToCard handles associating a label with a card.
// @Summary Add label to card
// @Description Associates an existing label with a specific card.
// @Tags Cards
// @Security ApiKeyAuth
// @Param cardID path string true "Card ID"
// @Param labelID path string true "Label ID"
// @Success 200 {object} models.Card "Card with label associated"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /cards/{cardID}/labels/{labelID} [post]
func AddLabelToCard(c *gin.Context) {
	cardID := c.Param("cardID")
	labelID := c.Param("labelID")

	card, err := cardService.AddLabelToCard(cardID, labelID)
	if err != nil {
		if strings.Contains(err.Error(), "card not found") || strings.Contains(err.Error(), "label not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		if strings.Contains(err.Error(), "label already associated") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to add label to card: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}

// RemoveLabelFromCard handles disassociating a label from a card.
// @Summary Remove label from card
// @Description Disassociates a label from a specific card.
// @Tags Cards
// @Security ApiKeyAuth
// @Param cardID path string true "Card ID"
// @Param labelID path string true "Label ID"
// @Success 200 {object} models.Card "Card with label disassociated"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /cards/{cardID}/labels/{labelID} [delete]
func RemoveLabelFromCard(c *gin.Context) {
	cardID := c.Param("cardID")
	labelID := c.Param("labelID")

	card, err := cardService.RemoveLabelFromCard(cardID, labelID)
	if err != nil {
		if strings.Contains(err.Error(), "card not found") || strings.Contains(err.Error(), "label not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		if strings.Contains(err.Error(), "label not associated") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to remove label from card: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, card)
}