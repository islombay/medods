package api

import (
	"github.com/gin-gonic/gin"
	"medods/api/handlers"
	"medods/pkg/logs"
	"medods/service"
)

func NewApiEndpoints(r *gin.RouterGroup, service service.ServiceInterface, log logs.LoggerInterface) {
	handler := handlers.New(log, service)

	auth := r.Group("/auth")
	{
		auth.POST("/login", handler.Login)
		auth.POST("/refresh", handler.Refresh)

		auth.POST("/register", handler.Register)
	}
}
