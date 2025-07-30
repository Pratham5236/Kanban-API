package controllers

import (
	"kanban-app/api/models"
	"kanban-app/api/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var projectService *services.ProjectService

func init() {
	projectService = services.NewProjectService()
}

// CreateProject handles creating a new project within an organization.
// @Summary Create a new project
// @Description Creates a new project within a specified organization. User must own the organization.
// @Tags Projects
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param project body models.CreateProjectRequest true "Project creation details"
// @Success 201 {object} models.Project "Project created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden (if user does not own organization)"
// @Failure 404 {object} models.ErrorResponse "Not Found (if organization not found)"
// @Failure 409 {object} models.ErrorResponse "Conflict (if project name already exists in organization)"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects [post]
func CreateProject(c *gin.Context) {
	userID, _ := c.Get("userID")
	orgID := c.Param("orgID")

	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	project, err := projectService.CreateProject(orgID, req.Name, req.Description, userID.(string))
	if err != nil {
		if strings.Contains(err.Error(), "project name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to create project: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// GetProjects handles retrieving all projects for a specific organization.
// @Summary Get all projects in an organization
// @Description Retrieves all projects within a specified organization. User must own the organization.
// @Tags Projects
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Produce json
// @Success 200 {array} models.Project "List of projects"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects [get]
func GetProjects(c *gin.Context) {
	orgID := c.Param("orgID")

	projects, err := projectService.GetProjectsByOrganizationID(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve projects: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetProjectByID handles retrieving a specific project within an organization.
// @Summary Get project by ID
// @Description Retrieves a specific project by its ID within a specified organization. User must own the organization.
// @Tags Projects
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Produce json
// @Success 200 {object} models.Project "Project details"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID} [get]
func GetProjectByID(c *gin.Context) {
	projectID := c.Param("projectID")

	project, err := projectService.GetProjectByID(projectID)
	if err != nil {
		if strings.Contains(err.Error(), "project not found") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to retrieve project: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject handles updating an existing project within an organization.
// @Summary Update a project
// @Description Updates a specific project by its ID within a specified organization. User must own the organization.
// @Tags Projects
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Param project body models.UpdateProjectRequest true "Project update details"
// @Success 200 {object} models.Project "Project updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 409 {object} models.ErrorResponse "Conflict"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID} [put]
func UpdateProject(c *gin.Context) {
	projectID := c.Param("projectID")

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
		return
	}

	updatedProject, err := projectService.UpdateProject(projectID, req)
	if err != nil {
		if strings.Contains(err.Error(), "project name already exists") {
			c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to update project: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProject)
}

// DeleteProject handles deleting a project within an organization.
// @Summary Delete a project
// @Description Deletes a specific project by its ID within a specified organization. User must own the organization.
// @Tags Projects
// @Security ApiKeyAuth
// @Param orgID path string true "Organization ID"
// @Param projectID path string true "Project ID"
// @Success 204 "No Content"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /organizations/{orgID}/projects/{projectID} [delete]
func DeleteProject(c *gin.Context) {
	projectID := c.Param("projectID")

	err := projectService.DeleteProject(projectID)
	if err != nil {
		if strings.Contains(err.Error(), "project not found or already deleted") {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "Failed to delete project: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}