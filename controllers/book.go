package controllers

import (
	"golang-rest-api-template/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}

// FindBooks godoc
// @Summary find books
// @Schemes
// @Description fetch all books data
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {string} book data
// @Router /books [get]
func FindBooks(c *gin.Context) {
	var books []models.Book
	models.DB.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": books})
}

// CreateBook godoc
// @Summary create book
// @Schemes
// @Description create book entry with title and author
// @Tags books
// @Accept json
// @Produce json
// @Success 201 {string} book data
// @Router /books [post]
func CreateBook(c *gin.Context) {
	var input models.CreateBook

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{Title: input.Title, Author: input.Author}

	models.DB.Create(&book)

	c.JSON(http.StatusCreated, gin.H{"data": book})
}

// FindBook godoc
// @Summary find book
// @Schemes
// @Description find book entry by id
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {string} book data
// @Router /books/{id} [get]
func FindBook(c *gin.Context) {
	var book models.Book

	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// UpdateBook godoc
// @Summary update book
// @Schemes
// @Description update book entry by id
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {string} book data
// @Router /books/{id} [put]
func UpdateBook(c *gin.Context) {
	var book models.Book
	var input models.UpdateBook

	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book not found"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&book).Updates(models.Book{Title: input.Title, Author: input.Author})

	c.JSON(http.StatusOK, gin.H{"data": book})
}

// DeleteBook godoc
// @Summary delete book
// @Schemes
// @Description delete book entry by id
// @Tags books
// @Accept json
// @Produce json
// @Success 204 {string} empty content
// @Router /books/{id} [delete]
func DeleteBook(c *gin.Context) {
	var book models.Book

	if err := models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "book not found"})
		return
	}

	models.DB.Delete(&book)

	c.JSON(http.StatusNoContent, gin.H{"data": true})
}
