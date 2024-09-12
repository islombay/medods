package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"medods/api/status"
	"medods/internal/error_list"
	"medods/pkg/logs"
	"medods/service"
)

type Handler struct {
	log     logs.LoggerInterface
	service service.ServiceInterface
}

func New(log logs.LoggerInterface, service service.ServiceInterface) *Handler {
	return &Handler{
		log: log, service: service,
	}
}

func (v *Handler) response(c *gin.Context, data status.Status) {
	if data.Code >= 200 && data.Code < 300 {
		c.JSON(data.Code, data)
		return
	}
	c.AbortWithStatusJSON(data.Code, data)
}

func (v *Handler) ParseError(err error) status.Status {
	if errors.Is(err, error_list.NotFound) {
		return status.StatusNotFound
	}

	return status.StatusInternal
}
