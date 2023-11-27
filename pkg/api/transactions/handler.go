package transactions

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

// FindTransaction godoc
// @Summary Find a transaction by ID
// @Description Get details of a transaction by its ID
// @Tags Transactions
// @Security JwtAuth
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} models.Transaction "Successfully retrieved transaction"
// @Failure 404 {string} string "Transaction not found"
// @Router /transactions/{id} [get]
func FindTransaction(c *gin.Context) {
	var transaction models.Transaction

	if err := database.DB.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// FindTransactions godoc
// @Summary Get all transactions with pagination
// @Description Get a list of all transactions with optional pagination
// @Tags Transactions
// @Security JwtAuth
// @Produce json
// @Param offset query int false "Offset for pagination" default(0)
// @Param limit query int false "Limit for pagination" default(10)
// @Success 200 {array} models.Transaction "Successfully retrieved list of transactions"
// @Router /transactions [get]
func FindTransactions(c *gin.Context) {
	var transactions []models.Transaction

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
	cacheKey := "transactions_offset_" + offsetQuery + "_limit_" + limitQuery

	// Try fetching the data from Redis first
	cachedTransactions, err := cache.Rdb.Get(cache.Ctx, cacheKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedTransactions), &transactions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal cached data"})
			return
		}
		c.JSON(http.StatusOK, transactions)
		return
	}

	// If cache missed, fetch data from the database
	database.DB.Offset(offset).Limit(limit).Find(&transactions)

	// Serialize transactions object and store it in Redis
	serializedTransactions, err := json.Marshal(transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
		return
	}
	err = cache.Rdb.Set(cache.Ctx, cacheKey, serializedTransactions, time.Minute).Err() // Here TTL is set to one hour
	if err != nil {
		log.Default().Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set cache"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction with the given input data
// @Tags Transactions
// @Security JwtAuth
// @Accept  json
// @Produce  json
// @Param   input     body   models.CreateTransaction   true   "Create transaction object"
// @Success 201 {object} models.Transaction "Successfully created transaction"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /transactions [post]
func CreateTransaction(c *gin.Context) {
	var input models.CreateTransaction

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, _ := time.Parse("2006-01-02", input.Date)

	transaction := models.Transaction{Account: input.Account, Date: date, Amount: input.Amount}

	database.DB.Create(&transaction)

	// Invalidate cache
	keysPattern := "transactions_offset_*"
	keys, err := cache.Rdb.Keys(cache.Ctx, keysPattern).Result()
	if err == nil {
		for _, key := range keys {
			cache.Rdb.Del(cache.Ctx, key)
		}
	}

	// Update account balance
	updateAccountBalance(uint(input.Account), input.Amount)

	c.JSON(http.StatusCreated, transaction)
}

// UpdateTransaction godoc
// @Summary Update a transaction by ID
// @Description Update the transaction details for the given ID
// @Tags Transactions
// @Security JwtAuth
// @Accept  json
// @Produce  json
// @Param id path string true "Transaction ID"
// @Param input body models.UpdateTransaction true "Update transaction object"
// @Success 200 {object} models.Transaction "Successfully updated transaction"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "transaction not found"
// @Router /transactions/{id} [put]
func UpdateTransaction(c *gin.Context) {
	var transaction models.Transaction
	var input models.UpdateTransaction

	if err := database.DB.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, _ := time.Parse("2006-01-02", input.Date)

	database.DB.Model(&transaction).Updates(models.Transaction{Account: input.Account, Date: date, Amount: input.Amount})

	// Update account balance
	updateAccountBalance(uint(input.Account), input.Amount)

	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction godoc
// @Summary Delete a transaction by ID
// @Description Delete the transaction with the given ID
// @Tags Transactions
// @Security JwtAuth
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 202 {object} models.Transaction "Successfully deleted transaction"
// @Failure 404 {string} string "transaction not found"
// @Router /transactions/{id} [delete]
func DeleteTransaction(c *gin.Context) {
	var transaction models.Transaction

	if err := database.DB.Where("id = ?", c.Param("id")).First(&transaction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	database.DB.Delete(&transaction)

	c.JSON(http.StatusAccepted, transaction)
}

// updateAccountBalance updates the balance of the given account
// by adding the given amount to the current balance
// Private function, not exposed to the API
func updateAccountBalance(account uint, amount float32) {
	var accountToUpdate models.Account

	database.DB.Where("account = ?", account).First(&accountToUpdate)
	accountToUpdate.Balance = accountToUpdate.Balance + amount
	database.DB.Save(&accountToUpdate)
}
