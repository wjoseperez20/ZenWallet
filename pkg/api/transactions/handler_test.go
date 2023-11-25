package transactions

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wjoseperez20/zenwallet/pkg/database"
	"github.com/wjoseperez20/zenwallet/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFindTransaction_SuccessfulRequest(t *testing.T) {
	// Given
	r := gin.Default()
	r.GET("/transactions/:id", FindTransaction)

	parseTime, err := time.Parse(time.RFC3339Nano, "2023-11-25T15:30:45.123456Z")
	require.NoError(t, err)

	dbMock, gormDB := setupTestDatabase(t)
	database.DB = gormDB
	mockTransaction := models.Transaction{ID: 1, Amount: 0.0, Date: parseTime, Account: 10001, CreatedAt: parseTime, UpdatedAt: parseTime}
	dbMock.ExpectQuery(`SELECT \* FROM "transactions" WHERE id = (.+) ORDER BY "transactions"."id" LIMIT 1`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "amount", "date", "account", "created_at", "updated_at"}).
			AddRow(mockTransaction.ID, mockTransaction.Amount, mockTransaction.Date, mockTransaction.Account, mockTransaction.CreatedAt, mockTransaction.UpdatedAt))

	// When
	w := performRequest(r, "GET", "/transactions/1")
	require.Equal(t, http.StatusOK, w.Code)

	var expected models.Transaction
	err = json.Unmarshal(w.Body.Bytes(), &expected)

	// Then
	require.NoError(t, err)
	require.Equal(t, mockTransaction.ID, expected.ID)

	// Verify all expectations were met
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindTransaction_NotFound(t *testing.T) {
	// Given
	r := gin.Default()
	r.GET("/transactions/:id", FindTransaction)

	dbMock, gormDB := setupTestDatabase(t)
	database.DB = gormDB
	dbMock.ExpectQuery(`SELECT \* FROM "transactions" WHERE id = (.+) ORDER BY "transactions"."id" LIMIT 1`).
		WithArgs("999").
		WillReturnError(gorm.ErrRecordNotFound)

	// When
	w := performRequest(r, "GET", "/transactions/999")
	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "transaction not found", errorResponse["error"])
}

// setupTestDatabase sets up a mock database for testing.
func setupTestDatabase(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	// Create a mock database for testing
	db, dbMock, err := sqlmock.New()
	require.NoError(t, err)

	// Replace the actual database with the mock database for testing
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return dbMock, gormDB
}

func performRequest(router *gin.Engine, method, path string, requestBody ...[]byte) *httptest.ResponseRecorder {
	var reqBody []byte
	if len(requestBody) > 0 {
		reqBody = requestBody[0]
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}
