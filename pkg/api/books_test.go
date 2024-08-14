package api

import (
	"bytes"
	"context"
	"encoding/json"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHealthcheck tests the Healthcheck handler.
func TestHealthcheck(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup the router
	r := gin.Default()
	r.GET("/", Healthcheck)

	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// Record the response
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, w.Code, "Expected response code to be 200")

	// Check the response body is what we expect.
	assert.Equal(t, "\"ok\"", w.Body.String(), "Expected response body to be 'ok'")
}

func TestFindBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a router
	r := gin.Default()

	// Setup middleware to inject app context
	r.Use(func(c *gin.Context) {
		ctx := context.Background()
		mockRedis := new(cache.MockRedisClient)
		mockDB := new(database.MockDB)

		// Setup the mock for Redis Get call
		key := "books_offset_0_limit_10"
		mockRedis.On("Get", ctx, key).Return(redis.NewStringResult("[]", nil)) // Adjust this as necessary

		appCtx := &AppContext{
			RedisClient: mockRedis,
			DB:          mockDB,
			Ctx:         &ctx,
		}
		c.Set("appCtx", appCtx)
		c.Next()
	})

	r.GET("/books", FindBooks)

	// Prepare and send the request
	req, _ := http.NewRequest(http.MethodGet, "/books?offset=0&limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code, "Expected response code to match")
	assert.JSONEq(t, `{"data":[]}`, w.Body.String(), "Expected response body to match JSON")
}

func TestCreateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		ctx := context.Background()
		mockRedis := new(cache.MockRedisClient)
		mockDB := new(database.MockDB)
		mockDB.On("Error").Return(nil) // Setup to return no error

		// Make sure to return the mockGormDB when Create is called
		mockDB.On("Create", mock.AnythingOfType("*models.Book")).Return(&gorm.DB{})

		// Mock Redis Keys call
		mockRedis.On("Keys", ctx, "books_offset_*").Return(redis.NewStringSliceResult([]string{}, nil))

		// Mock Redis Del call (in case it is used after Keys)
		mockRedis.On("Del", ctx, mock.AnythingOfType("string")).Return(redis.NewIntResult(1, nil))

		appCtx := &AppContext{
			RedisClient: mockRedis,
			DB:          mockDB,
			Ctx:         &ctx,
		}
		c.Set("appCtx", appCtx)
		c.Next()
	})

	r.POST("/books", CreateBook)

	testBook := models.CreateBook{
		Title:  "Test Book",
		Author: "Test Author",
	}
	body, _ := json.Marshal(testBook)
	req := httptest.NewRequest("POST", "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	expectedJSON := `{
		"data": {
			"title": "Test Book",
			"author": "Test Author",
			"id": 0,
			"created_at": "0001-01-01T00:00:00Z",
			"updated_at": "0001-01-01T00:00:00Z"
		}
	}`
	assert.JSONEq(t, expectedJSON, w.Body.String(), "Expected response body to match JSON")
}
