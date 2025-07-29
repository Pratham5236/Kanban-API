package controllers

import (
	"kanban-app/api/models"
	"kanban-app/api/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var organizationService *services.OrganizationService

func init() {
	organizationService = services.NewOrganizationService()
}

// TODO check if i really even need this, since its bbeing handled by the middleware
func getUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "User ID not found in context"})
		return "", false
	}
	return userID.(string), true
}

// CreateOrganization handles creating a new organization.
// @Summary Create a new organization
// @Description Creates a new organization associated with the authenticated user.
// @Tags Organizations
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param organization body models.CreateOrganizationRequest true "Organization creation details"
// @Success 201 {object} models.Organization "Organization created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request (e.g., validation error)"
// @Failure 401 {object} models.ErrorResponse "Unauthorized (if token is missing or invalid)"
// @Failure 409 {object} models.ErrorResponse "Conflict (if organization name already exists)"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations [post]
func CreateOrganization(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	var req models.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	org, err := organizationService.CreateOrganization(req.Name, userID)
	if err != nil {
		if strings.Contains(err.Error(), "organization name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to create organization: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// GetOrganizations handles retrieving all organizations for the authenticated user.
// @Summary Get user's organizations
// @Description Retrieves all organizations owned by the authenticated user.
// @Tags Organizations
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} models.Organization "List of organizations"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations [get]
func GetOrganizations(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orgs, err := organizationService.GetOrganizationsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve organizations: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, orgs)
}

// GetOrganizationByID handles retrieving a single organization by its ID.
// @Summary Get organization by ID
// @Description Retrieves a specific organization by its ID, ensuring the authenticated user is the owner.
// @Tags Organizations
// @Security ApiKeyAuth
// @Param id path string true "Organization ID"
// @Produce json
// @Success 200 {object} models.Organization "Organization details"
// @Failure 400 {object} models.ErrorResponse "Bad Request (e.g., invalid ID format)"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden (if user is not the owner)"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{id} [get]
func GetOrganizationByID(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orgID := c.Param("id")

	org, err := organizationService.GetOrganizationByID(orgID)
	if err != nil {
		if strings.Contains(err.Error(), "organization not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retreive organization: " + err.Error()})
		return
	}

	if org.OwnerID != userID {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "You do not have permission to access this organization"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// UpdateOrganization handles updating an existing organization.
// @Summary Update an organization
// @Description Updates a specific organization by its ID, ensuring the authenticated user is the owner.
// @Tags Organizations
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Param organization body models.UpdateOrganizationRequest true "Organization update details"
// @Success 200 {object} models.Organization "Organization updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 409 {object} models.ErrorResponse "Conflict (if new name already exists)"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{id} [put]
func UpdateOrganization(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orgID := c.Param("id")

	existingOrg, err := organizationService.GetOrganizationByID(orgID)
	if err != nil {
		if strings.Contains(err.Error(), "organization not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error getting organization: " + err.Error()})
		return
	}

	if existingOrg.OwnerID != userID {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "You do not have permission to update this organization"})
		return
	}

	var req models.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	updatedOrg, err := organizationService.UpdateOrganization(orgID, req)
	if err != nil {
		if strings.Contains(err.Error(), "organization name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to update organization: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedOrg)
}

// DeleteOrganization handles deleting an organization.
// @Summary Delete an organization
// @Description Deletes a specific organization by its ID, ensuring the authenticated user is the owner.
// @Tags Organizations
// @Security ApiKeyAuth
// @Param id path string true "Organization ID"
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{id} [delete]
func DeleteOrganization(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		return
	}

	orgID := c.Param("id")

	existingOrg, err := organizationService.GetOrganizationByID(orgID)
	if err != nil {
		if strings.Contains(err.Error(), "organization not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Error getting organization: " + err.Error()})
		return
	}

	if existingOrg.OwnerID != userID {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "You do not have permission to update this organization"})
		return
	}

	err = organizationService.DeleteOrganization(orgID)
	if err != nil {
		if strings.Contains(err.Error(), "organization not found or already deleted") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to delete organization: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
