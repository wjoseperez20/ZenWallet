package files

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// FindFile godoc
// @Summary Find file by ID
// @Description Get details of a file by its ID
// @Tags Files
// @Security JwtAuth
// @Accept json
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} models.File "Successfully retrieved file"
// @Failure 404 {string} string "File not found"
// @Router /files/{id} [get]
func FindFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "FindFile"})
}

// FindFiles godoc
// @Summary Get all files with pagination
// @Description Get a list of all files with optional pagination
// @Tags Files
// @Security JwtAuth
// @Accept json
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} models.File "Successfully retrieved list of files"
// @Router /files [get]
func FindFiles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "FindFiles"})
}

// UploadFile godoc
// @Summary Upload a new file
// @Description
// @Tags Files
// @Security JwtAuth
// @Accept  json
// @Produce  json
// @Param   input     body   models.UploadFile   true   "Upload file"
// @Success 201 {object} models.File "Successfully created file"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /files [post]
func UploadFile(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "UploadFile"})
}
