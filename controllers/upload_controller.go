package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadFile handles file uploads.
// @Summary Upload a file
// @Description Uploads a single file to the server.
// @Tags Uploads
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string "File uploaded successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /uploads [post]
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("get form file error: %s", err.Error())})
		return
	}

	// Create the uploads directory if it doesn't exist
	uploadsDir := "./uploads"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		os.Mkdir(uploadsDir, 0755)
	}

	filename := filepath.Base(file.Filename)
	dst := filepath.Join(uploadsDir, filename)

	// Save the uploaded file to the specified destination
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("upload file error: %s", err.Error())})
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "file_url": fileURL})
}
