package files

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/zenwallet/pkg/amazon"
	"github.com/wjoseperez20/zenwallet/pkg/cache"
	"github.com/wjoseperez20/zenwallet/pkg/database"
	"github.com/wjoseperez20/zenwallet/pkg/models"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	s3Bucket       = "zenwallet-bucket"
	s3BucketFolder = "files"
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
// @Router /files/upload [post]
func UploadFile(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// Create an S3 client
	svc := s3.New(amazon.Aws)

	// Generate a unique file name
	fileName := header.Filename

	// Upload the file to S3
	err = uploadToS3(svc, file, fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a file object
	fileObj := models.File{Name: fileName, Location: "S3"}

	// Save the file object to the database
	database.DB.Create(&fileObj)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

// ProcessFile godoc
// @Summary Process a file
// @Description
// @Tags Files
// @Security JwtAuth
// @Accept  json
// @Produce  json
// @Param   input     body   models.ProcessFile   true   "Upload file"
// @Success 201 {object} models.File "Successfully processed file"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /files/process [post]
func ProcessFile(c *gin.Context) {
	var input models.ProcessFile

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a file object
	file := models.File{Name: input.Name}

	// Check if the file exists in the database
	if err := database.DB.Where("name = ?", file.Name).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// Check if the file has already been processed
	if file.Processed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file has already been processed"})
		return
	}

	// Create an S3 client
	svc := s3.New(amazon.Aws)

	// Download the file to S3
	err := downloadFromS3(svc, file.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process the file
	err = insertTransactions(file.Name)
	if err != nil {
		database.DB.Model(&file).Updates(models.File{Output: err.Error()})

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Delete the file from the server
	err = os.Remove(file.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the file object in the database
	database.DB.Model(&file).Updates(models.File{Processed: true})

	c.JSON(http.StatusOK, gin.H{"message": "File processed successfully"})
}

// uploadToS3 godoc
// @Summary Upload a file
// Private function to upload a file to S3
func uploadToS3(svc *s3.S3, file multipart.File, fileName string) error {
	// Set up the S3 upload parameters
	params := &s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", s3BucketFolder, fileName)),
		Body:   file,
	}

	// Perform the upload
	_, err := svc.PutObject(params)
	if err != nil {
		return err
	}

	return err
}

// downloadFile godoc
// @Summary Download a file
// Private function to download a file from S3
func downloadFromS3(svc *s3.S3, fileName string) error {
	// Set up the S3 download parameters
	params := &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", s3BucketFolder, fileName)),
	}

	// Create a file to write the S3 object content to
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Perform the download
	result, err := svc.GetObject(params)
	if err != nil {
		return err
	}

	// Copy the S3 object contents to the file
	_, err = io.Copy(file, result.Body)
	if err != nil {
		return err
	}

	return nil
}

// processFile godoc
// @Summary Process a file
// Private function to process a file
func insertTransactions(fileName string) error {
	// process csv file
	csvTransactions, err := readCSV(fileName)
	if err != nil {
		return err
	}

	// save transactions to database
	for _, input := range csvTransactions {

		// check if account exists
		var account models.Account
		if err := database.DB.Where("account = ?", input.Account).First(&account).Error; err != nil {
			return fmt.Errorf("account %d not found", input.Account)
		}

		// Create a transaction object
		transaction := models.Transaction{Account: input.Account, Date: input.Date, Amount: input.Amount}

		// Save the transaction object to the database
		database.DB.Create(&transaction)

		// Update account balance
		updateAccountBalance(uint(transaction.Account), transaction.Amount)
	}

	return nil
}

// readCSV godoc
// @Summary Read a CSV file
// Private function to read a CSV file
func readCSV(filename string) ([]models.Transaction, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and ignore the header line
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction

	for _, record := range records {
		account, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return nil, err
		}

		date, err := time.Parse("2006-01-02", record[2])
		if err != nil {
			return nil, err
		}

		amount, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}

		transaction := models.Transaction{
			Account: int(account),
			Date:    date,
			Amount:  float32(amount),
		}

		transactions = append(transactions, transaction)

	}

	return transactions, nil
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
