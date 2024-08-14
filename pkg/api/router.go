package api

import (
	"context"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/middleware"
	"time"

	docs "golang-rest-api-template/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"golang.org/x/time/rate"
)

// AppContext holds shared resources like database and Redis client
type AppContext struct {
	DB          database.DBInterface
	RedisClient cache.CacheInterface
	Ctx         *context.Context
}

// NewAppContext creates a new AppContext
func NewAppContext(db database.DBInterface, redisClient cache.CacheInterface, ctx *context.Context) *AppContext {
	return &AppContext{
		DB:          db,
		RedisClient: redisClient,
		Ctx:         ctx,
	}
}

func ContextMiddleware(appCtx *AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCtx", appCtx)
		c.Next()
	}
}

func NewRouter(logger *zap.Logger, mongoCollection *mongo.Collection, db database.DBInterface, redisClient cache.CacheInterface, ctx *context.Context) *gin.Engine {
	appCtx := NewAppContext(db, redisClient, ctx)

	r := gin.Default()
	r.Use(ContextMiddleware(appCtx))

	//r.Use(gin.Logger())
	r.Use(middleware.Logger(logger, mongoCollection))
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/", Healthcheck)
		v1.GET("/books", middleware.APIKeyAuth(), FindBooks)
		v1.POST("/books", middleware.APIKeyAuth(), middleware.JWTAuth(), CreateBook)
		v1.GET("/books/:id", middleware.APIKeyAuth(), FindBook)
		v1.PUT("/books/:id", middleware.APIKeyAuth(), UpdateBook)
		v1.DELETE("/books/:id", middleware.APIKeyAuth(), DeleteBook)

		v1.POST("/login", middleware.APIKeyAuth(), LoginHandler)
		v1.POST("/register", middleware.APIKeyAuth(), RegisterHandler)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
