package main

import (
	"golang-rest-api-template/auth"
	"golang-rest-api-template/cache"
	"golang-rest-api-template/controllers"
	"golang-rest-api-template/middleware"
	"golang-rest-api-template/models"
	"log"
	"time"

	docs "golang-rest-api-template/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8001
// @BasePath  /api/v1

// @securityDefinitions.apikey JwtAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cache.InitRedis()

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r.Use(gin.Logger())
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.SecurityMiddleware())
		r.Use(middleware.XssMiddleware())
	}
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RateLimitMiddleware(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	models.ConnectDatabase()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", controllers.Healthcheck)
		v1.GET("/books", middleware.APIKeyAuthMiddleware(), controllers.FindBooks)
		v1.POST("/books", middleware.APIKeyAuthMiddleware(), middleware.AuthenticateJWT(), controllers.CreateBook)
		v1.GET("/books/:id", middleware.APIKeyAuthMiddleware(), controllers.FindBook)
		v1.PUT("/books/:id", middleware.APIKeyAuthMiddleware(), controllers.UpdateBook)
		v1.DELETE("/books/:id", middleware.APIKeyAuthMiddleware(), controllers.DeleteBook)

		v1.POST("/login", middleware.APIKeyAuthMiddleware(), auth.LoginHandler)
		v1.POST("/register", middleware.APIKeyAuthMiddleware(), auth.RegisterHandler)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
