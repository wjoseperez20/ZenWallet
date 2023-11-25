package files

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/zenwallet/pkg/cache"
	"github.com/wjoseperez20/zenwallet/pkg/database"
	"github.com/wjoseperez20/zenwallet/pkg/models"
	"log"
	"net/http"
	"strconv"
	"time"
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
	var file models.File

	if err := database.DB.Where("id = ?", c.Param("id")).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, file)
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
	var files []models.File

	// Get query params
	offsetQuery := c.DefaultQuery("offset", "0")
	limitQuery := c.DefaultQuery("limit", "10")

	// Convert query params to integers
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset format"})
		return
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}

	// Create a cache key based on query params
	cacheKey := "files_offset_" + offsetQuery + "_limit_" + limitQuery

	// Try fetching the data from Redis first
	cachedFiles, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedFiles), &files)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal cached data"})
			return
		}
		c.JSON(http.StatusOK, files)
		return
	}

	// If cache missed, fetch data from the database
	database.DB.Offset(offset).Limit(limit).Find(&files)

	// Serialize files object and store it in Redis
	serializedFiles, err := json.Marshal(files)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, cacheKey, serializedFiles, time.Minute).Err() // Here TTL is set to one hour
	if err != nil {
		log.Default().Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache"})
		return
	}

	c.JSON(http.StatusOK, files)
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
	var input models.UploadFile

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file := models.File{Name: input.Name}

	database.DB.Create(&file)

	// Invalidate cache
	keysPattern := "files_offset_*"
	keys, err := cache.Rdb.Keys(cache.Ctx, keysPattern).Result()
	if err == nil {
		for _, key := range keys {
			cache.Rdb.Del(cache.Ctx, key)
		}
	}

	c.JSON(http.StatusCreated, file)
}
