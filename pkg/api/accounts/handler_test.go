package accounts

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
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

func TestFindAccount_SuccessfulRequest(t *testing.T) {
	// Given
	r := gin.Default()
	r.GET("/accounts/:account", FindAccount)

	parseTime, err := time.Parse(time.RFC3339Nano, "2023-11-25T15:30:45.123456Z")
	require.NoError(t, err)

	dbMock, gormDB := setupTestDatabase(t)
	database.DB = gormDB
	mockAccount := models.Account{ID: 10001, Client: "test", Email: "test@email.com", Account: 10001, Balance: 1.0, CreatedAt: parseTime, UpdatedAt: parseTime}
	dbMock.ExpectQuery(`SELECT \* FROM "accounts" WHERE account = (.+) ORDER BY "accounts"."account" LIMIT 1`).
		WithArgs("10001").
		WillReturnRows(sqlmock.NewRows([]string{"id", "client", "email", "account", "balance", "created_at", "updated_at"}).
			AddRow(mockAccount.ID, mockAccount.Client, mockAccount.Email, mockAccount.Account, mockAccount.Balance, mockAccount.CreatedAt, mockAccount.UpdatedAt))

	// When
	w := performRequest(r, "GET", "/accounts/10001")
	require.Equal(t, http.StatusOK, w.Code)

	var expected models.Account
	err = json.Unmarshal(w.Body.Bytes(), &expected)

	// Then
	require.NoError(t, err)
	require.Equal(t, mockAccount.ID, expected.ID)

	// Verify all expectations were met
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindAccount_NotFound(t *testing.T) {
	// Given
	r := gin.Default()
	r.GET("/accounts/:account", FindAccount)

	dbMock, gormDB := setupTestDatabase(t)
	database.DB = gormDB
	dbMock.ExpectQuery(`SELECT \* FROM "accounts" WHERE account = (.+) ORDER BY "accounts"."account" LIMIT 1`).
		WithArgs("999").
		WillReturnError(gorm.ErrRecordNotFound)

	// When
	w := performRequest(r, "GET", "/accounts/999")
	require.Equal(t, http.StatusNotFound, w.Code)

	// Verify all expectations were met
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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

// performRequest performs an HTTP request and returns the response recorder.
func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}
