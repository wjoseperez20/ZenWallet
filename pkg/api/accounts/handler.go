package accounts

import (
	"encoding/json"
	"github.com/wjoseperez20/zenwallet/pkg/cache"
	"github.com/wjoseperez20/zenwallet/pkg/database"
	"github.com/wjoseperez20/zenwallet/pkg/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// FindAccount godoc
// @Summary Find an account by ID
// @Description Get details of an account by its ID
// @Tags Accounts
// @Security JwtAuth
// @Produce json
// @Param id path string true "Account ID"
// @Success 200 {object} models.Account "Successfully retrieved account"
// @Failure 404 {string} string "Account not found"
// @Router /accounts/{id} [get]
func FindAccount(c *gin.Context) {
	var account models.Account

	if err := database.DB.Where("account = ?", c.Param("account")).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// FindAccounts godoc
// @Summary Get all accounts with pagination
// @Description Get a list of all accounts with optional pagination
// @Tags Accounts
// @Security JwtAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} models.Account "Successfully retrieved list of accounts"
// @Router /accounts [get]
func FindAccounts(c *gin.Context) {
	var accounts []models.Account

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
	cacheKey := "accounts_offset_" + offsetQuery + "_limit_" + limitQuery

	// Try fetching the data from Redis first
	cachedAccounts, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedAccounts), &accounts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal cached data"})
			return
		}
		c.JSON(http.StatusOK, accounts)
		return
	}

	// If cache missed, fetch data from the database
	database.DB.Offset(offset).Limit(limit).Find(&accounts)

	// Serialize accounts object and store it in Redis
	serializedAccounts, err := json.Marshal(accounts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, cacheKey, serializedAccounts, time.Minute).Err() // Here TTL is set to one hour
	if err != nil {
		log.Default().Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account with the given input data
// @Tags Accounts
// @Security JwtAuth
// @Accept  json
// @Produce  json
// @Param   input     body   models.CreateAccount   true   "Create account object"
// @Success 201 {object} models.Account "Successfully created account"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /accounts [post]
func CreateAccount(c *gin.Context) {
	var input models.CreateAccount

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account := models.Account{Client: input.Client, Email: input.Email}

	database.DB.Create(&account)

	// Invalidate cache
	keysPattern := "accounts_offset_*"
	keys, err := cache.Rdb.Keys(cache.Ctx, keysPattern).Result()
	if err == nil {
		for _, key := range keys {
			cache.Rdb.Del(cache.Ctx, key)
		}
	}

	c.JSON(http.StatusCreated, account)
}

// UpdateAccount godoc
// @Summary Update an account by ID
// @Description Update the account details for the given ID
// @Tags Accounts
// @Security JwtAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Account ID"
// @Param input body models.UpdateAccount true "Update account object"
// @Success 200 {object} models.Account "Successfully updated account"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "account not found"
// @Router /accounts/{id} [put]
func UpdateAccount(c *gin.Context) {
	var account models.Account
	var input models.UpdateAccount

	if err := database.DB.Where("account = ?", c.Param("account")).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&account).Updates(models.Account{Client: input.Client, Email: input.Email})

	c.JSON(http.StatusOK, account)
}

// DeleteAccount godoc
// @Summary Delete an account by ID
// @Description Delete the account with the given ID
// @Tags Accounts
// @Security JwtAuth
// @Produce json
// @Param id path string true "Account ID"
// @Success 202 {object} models.Account "Successfully deleted account"
// @Failure 404 {string} string "account not found"
// @Router /accounts/{id} [delete]
func DeleteAccount(c *gin.Context) {
	var account models.Account

	if err := database.DB.Where("account = ?", c.Param("account")).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
		return
	}

	database.DB.Delete(&account)

	c.JSON(http.StatusAccepted, account)
}
