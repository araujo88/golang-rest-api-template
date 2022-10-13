package main

import (
	"golang-rest-api-template/controllers"
	"golang-rest-api-template/models"
	"log"

	docs "golang-rest-api-template/docs"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	if gin.Mode() == gin.ReleaseMode {
		r.Use(secure.New(secure.Config{
			//AllowedHosts:          []string{"example.com", "ssl.example.com"},
			//SSLRedirect:           true,
			//SSLHost:               "ssl.example.com",
			STSSeconds:            315360000,
			STSIncludeSubdomains:  true,
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
			IENoOpen:              true,
			ReferrerPolicy:        "strict-origin-when-cross-origin",
			SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		}))
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://127.0.0.1",
			"http://127.0.0.1:8001",
			"http://localhost",
			"http://localhost:8001"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		//AllowOriginFunc: func(origin string) bool {
		//	return origin == "https://github.com"
		//},
		MaxAge: 12 * time.Hour,
	}))

	models.ConnectDatabase()

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/helloworld", controllers.Helloworld)
		v1.GET("/books", controllers.FindBooks)
		v1.POST("/books", controllers.CreateBook)
		v1.GET("/books/:id", controllers.FindBook)
		v1.PUT("/books/:id", controllers.UpdateBook)
		v1.DELETE("/books/:id", controllers.DeleteBook)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
