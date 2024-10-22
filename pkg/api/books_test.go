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

func TestCreateBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	mockCache := cache.NewMockCache(ctrl)
	ctx := context.Background()

	repo := NewBookRepository(mockDB, mockCache, &ctx)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/books", func(c *gin.Context) {
		// Set the appCtx in the Gin context
		c.Set("appCtx", repo)
		repo.CreateBook(c)
	})

	// Example data for the test
	inputBook := models.CreateBook{Title: "New Book", Author: "New Author"}
	requestBody, err := json.Marshal(inputBook)
	if err != nil {
		t.Fatalf("Failed to marshal input book data: %v", err)
	}

	// Set up database mock to simulate successful book creation
	mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(book *models.Book) *gorm.DB {
		// Normally, you might simulate setting an ID or other fields modified by the DB
		return &gorm.DB{Error: nil}
	})

	// Set up cache mock to simulate key retrieval and deletion
	keyPattern := "books_offset_*"
	mockCache.EXPECT().Keys(ctx, keyPattern).Return(redis.NewStringSliceResult([]string{"books_offset_0_limit_10"}, nil))
	mockCache.EXPECT().Del(ctx, "books_offset_0_limit_10").Return(redis.NewIntResult(1, nil))

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to create the HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Serve the HTTP request
	r.ServeHTTP(w, req)

	// Assertions to check the response
	assert.Equal(t, http.StatusCreated, w.Code, "Expected HTTP status code 201")
	assert.Contains(t, w.Body.String(), "New Book", "Response body should contain the book title")
}

func TestFindBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := database.NewMockDatabase(ctrl)
	ctx := context.Background()
	repo := NewBookRepository(mockDB, nil, &ctx)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/book/:id", repo.FindBook)

	// Prepare test data
	expectedBook := models.Book{
		ID:     1,
		Title:  "Effective Go",
		Author: "Robert Griesemer",
	}

	// Mock expectations

	// Mock the Where method
	mockDB.EXPECT().
		Where("id = ?", "1").
		DoAndReturn(func(query interface{}, args ...interface{}) database.Database {
			// Return mockDB to allow method chaining
			return mockDB
		}).Times(1)

	// Mock the First method
	mockDB.EXPECT().
		First(gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) database.Database {
			if b, ok := dest.(*models.Book); ok {
				*b = expectedBook
			}
			return mockDB
		}).Times(1)

	// Mock the Error method or field access
	mockDB.EXPECT().
		Error().
		Return(nil).
		Times(1)

	// Perform the request
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/book/1", nil)
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    models.Book `json:"data"`
	}

	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedBook.ID, response.Data.ID)
	assert.Equal(t, expectedBook.Title, response.Data.Title)
	assert.Equal(t, expectedBook.Author, response.Data.Author)
}

func TestDeleteBook(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock for the database
	mockDB := database.NewMockDatabase(ctrl)
	ctx := context.Background()
	repo := NewBookRepository(mockDB, nil, &ctx)

	// Set up Gin for testing
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/book/:id", repo.DeleteBook)

	// Prepare the book data
	existingBook := models.Book{
		ID:     1,
		Title:  "Test Book",
		Author: "Test Author",
	}

	// Mock Where to return the existingBook for chaining
	mockDB.EXPECT().
		Where("id = ?", "1").
		Return(mockDB).Times(1)

	// Mock First to load the existingBook and return mockDB
	mockDB.EXPECT().
		First(gomock.Any()).
		DoAndReturn(func(dest interface{}, conds ...interface{}) database.Database {
			if b, ok := dest.(*models.Book); ok {
				*b = existingBook
			}
			return mockDB
		}).Times(1)

	// Mock Delete method
	mockDB.EXPECT().
		Delete(&existingBook).
		Return(&gorm.DB{Error: nil}).Times(1)

	// Mock Error method to return nil
	mockDB.EXPECT().Error().Return(nil).AnyTimes()

	// Perform the DELETE request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/book/1", nil)
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusNoContent, w.Code)
}
