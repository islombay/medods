package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"medods/api/status"
	"medods/internal/error_list"
	"medods/pkg/logs"
	"medods/service"
	"time"
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
	} else if errors.Is(err, error_list.Unauthorized) {
		return status.StatusUnauthorized
	} else if errors.Is(err, error_list.Forbidden) {
		return status.StatusForbidden
	}

	return status.StatusInternal
}

func GetValue[T any](c *gin.Context, key string) (T, bool) {
	var zeroValue T
	val, ok := c.Get(key)
	if !ok {
		return zeroValue, false
	}

	assertedValue, ok := val.(T)
	if !ok {
		return zeroValue, false
	}

	return assertedValue, true
}

func (v *Handler) NewContext(c *gin.Context, sec int, keys ...string) (context.Context, context.CancelFunc) {
	ctx := context.Background()
	for _, key := range keys {
		if val, ok := c.Get(key); ok {
			ctx = context.WithValue(ctx, key, val)
		}
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(sec))

	return ctx, cancel
}
