package api

import (
	"context"
	"encoding/json"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewBookRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	mockCtx := context.Background()

	repo := NewBookRepository(mockDB, mockCache, &mockCtx)

	assert.NotNil(t, repo, "NewBookRepository should return a non-nil instance of bookRepository")
	assert.Equal(t, mockDB, repo.DB, "DB should be set to the mock database instance")
	assert.Equal(t, mockCache, repo.RedisClient, "RedisClient should be set to the mock cache instance")
}

func TestHealthcheck(t *testing.T) {
	// Set up the mock controller and the mocked dependencies
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Set up the Gin context with a response recorder
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	_, router := gin.CreateTestContext(recorder)

	// Create a mock repository and expect the Healthcheck method to be called
	mockRepo := NewMockBookRepository(ctrl)
	mockRepo.EXPECT().Healthcheck(gomock.Any()).Do(func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok") // Explicitly setting the response here
	})

	// Setting up a basic GET route to test Healthcheck
	router.GET("/healthcheck", mockRepo.Healthcheck)

	// Perform the GET request
	req, _ := http.NewRequest(http.MethodGet, "/healthcheck", nil)
	router.ServeHTTP(recorder, req)

	// Check the response
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "\"ok\"", recorder.Body.String())
}

func TestFindBooks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	mockGormDB := database.NewMockDatabase(ctrl) // Correct type for GORM DB operations
	ctx := context.Background()

	repo := NewBookRepository(mockDB, mockCache, &ctx)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/books", repo.FindBooks)

	// Set up common mock expectations
	mockGormDB.EXPECT().Find(gomock.Any()).DoAndReturn(func(books *[]models.Book) *gorm.DB {
		*books = append(*books, models.Book{Title: "New Book", Author: "New Author"})
		return &gorm.DB{Error: nil} // Assume this is the struct provided by the actual Gorm package
	}).AnyTimes()

	books := []models.Book{{Title: "Book One", Author: "Author One"}}
	cachedData, _ := json.Marshal(books)
	mockCache.EXPECT().Get(ctx, "books_offset_0_limit_10").Return(redis.NewStringResult(string(cachedData), nil))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books?offset=0&limit=10", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Book One")
}
