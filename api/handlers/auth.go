package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"medods/api/status"
	"medods/internal/model"
	"medods/pkg/helper"
	"medods/pkg/logs"
	"time"
)

// Login
// @id login
// @router /api/auth/login [post]
// @summary Login
// @description Login Get Access/ Refresh token
// @tags auth
// @param product body model.LoginRequest true "Request"
// @success 200 {object} model.TokenPair "Tokens"
// @failure 500 {object} status.Status "Internal server error"
func (v *Handler) Login(c *gin.Context) {
	var m model.LoginRequest
	if err := c.BindJSON(&m); err != nil {
		v.response(c, status.StatusBadRequest)
		return
	}

	// when uuid is invalid
	if !helper.IsUUID(m.UserId) {
		v.response(c, status.StatusBadRequest.AddError("uuid", "invalid"))
		return
	}

	// TODO: handle cases when the IP is inside header (ex. Nginx)
	m.IP = c.RemoteIP()

	v.log.Debug("IP",
		logs.String("c.RemoteIP()", c.RemoteIP()),
		logs.String("c.ClientIP()", c.ClientIP()),
		logs.String("c.Request.RemoteAddr", c.Request.RemoteAddr))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tokenPair, err := v.service.Auth().Login(ctx, m)
	if err != nil {
		v.response(c, v.ParseError(err))
		return
	}

	v.response(c, status.StatusOk.AddData(tokenPair))
}

func (v *Handler) Refresh(c *gin.Context) {

}

// Register
// @id register
// @router /api/auth/register [post]
// @summary Register new user
// @description Functionality to test Login function
// @tags auth
// @param register body model.Register true "Request"
// @success 200 {object} model.TokenPair "Tokens"
// @failure 500 {object} status.Status "Internal server error"
func (v *Handler) Register(c *gin.Context) {
	var m model.Register
	if err := c.BindJSON(&m); err != nil {
		v.response(c, status.StatusBadRequest)
		return
	}

	m.IP = helper.ParseIP(c)

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	//defer cancel()
	ctx := context.Background()

	tokenPair, err := v.service.Auth().Register(ctx, m)
	if err != nil {
		v.response(c, v.ParseError(err))
		return
	}
	v.response(c, status.StatusOk.AddData(tokenPair))
}
