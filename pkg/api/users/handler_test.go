package users

import (
	"bytes"
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

func TestLoginUser(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/login", LoginUser)

	parseTime, err := time.Parse(time.RFC3339Nano, "2023-11-25T15:30:45.123456Z")
	incomingUser := models.User{
		Username: "test",
		Password: "test",
	}

	dbMock, gormDB := setupTestDatabase(t)
	database.DB = gormDB
	mockUser := models.User{Username: "test", Password: "$2a$14$7z17lzN8ckCiGEQQdbQ2c.XsnJYDunu8SQ1H9BG9EqT4FpVwez68K", CreatedAt: parseTime, UpdatedAt: parseTime}
	dbMock.ExpectQuery(`SELECT \* FROM "users" WHERE username = (.+) ORDER BY "users"."username" LIMIT 1`).
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"username", "password", "created_at", "updated_at"}).
			AddRow(mockUser.Username, mockUser.Password, mockUser.CreatedAt, mockUser.UpdatedAt))

	// When
	w := performRequest(r, "POST", "/login", toJSON(incomingUser))
	require.Equal(t, http.StatusOK, w.Code)

	var expected map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &expected)

	// Then
	require.NoError(t, err)
	require.NotNil(t, expected["token"])
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

func toJSON(v interface{}) []byte {
	result, _ := json.Marshal(v)
	return result
}
