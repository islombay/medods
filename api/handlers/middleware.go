package handlers

import (
	"github.com/gin-gonic/gin"
	"medods/api/status"
	"medods/internal/model"
	"medods/pkg/jwt"
	"medods/pkg/logs"
)

type AuthorizationHeaders struct {
	Token string `header:"Authorization" binding:"required"`
}

func (v *Handler) haveToken(c *gin.Context) (*jwt.AccessTokenClaims, *status.Status) {
	var m AuthorizationHeaders
	if err := c.ShouldBindHeader(&m); err != nil {
		return nil, &status.StatusUnauthorized
	}

	claims, err := jwt.ParseToken(m.Token)
	if err != nil {
		return nil, &status.StatusUnauthorized
	}

	return claims, nil
}

func (v *Handler) MiddlewareIsUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, statusResponse := v.haveToken(c)
		if statusResponse != nil {
			v.response(c, *statusResponse)
			v.log.Debug("exiting here")
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func (v *Handler) MiddlewareDeviceTrack() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m model.DeviceInfo
		m.IP = c.ClientIP()
		m.UserAgent = c.Request.UserAgent()

		v.log.Debug("tracking device info", logs.Any("device_info", m))

		c.Set("device_info", m)
		c.Next()
	}
}
