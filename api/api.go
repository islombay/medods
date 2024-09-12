package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "medods/api/docs"
	"medods/config"
	"medods/pkg/logs"
	"medods/service"
)

func New(r *gin.Engine, service service.ServiceInterface, cfg config.Config, log logs.LoggerInterface) {
	r.Use(customCORSMiddleware())

	// @title Medods Auth Service
	// @description Authorization Service for Medods
	// @version 1.0

	// @securityDefinitions.apikey ApiKeyAuth
	// @in header
	// @name Authorization

	// @host localhost:8095
	api := r.Group("/api")

	NewApiEndpoints(api, service, log)

	r.GET("/sw/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
	))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"ping": "pong"})
	})
}

func customCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE, HEAD")
		c.Header("Access-Control-Allow-Headers", "Platform-Id, Content-Type, Content-Length, Accept-Encoding, X-CSF-TOKEN, Authorization, Cache-Control")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
