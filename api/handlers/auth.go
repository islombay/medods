package handlers

import (
	"github.com/gin-gonic/gin"
	"medods/api/status"
	"medods/internal/model"
	"medods/pkg/helper"
)

// Login
// @id login
// @router /api/auth/login [post]
// @summary Login
// @description Login Get Access/ Refresh token
// @tags auth
// @param request body model.LoginRequest true "Request"
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

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	ctx, cancel := v.NewContext(c, 5, "device_info")
	defer cancel()

	tokenPair, err := v.service.Auth().Login(ctx, m)
	if err != nil {
		v.response(c, v.ParseError(err))
		return
	}

	v.response(c, status.StatusOk.AddData(tokenPair))
}

// Refresh
// @id refresh
// @router /api/auth/refresh [post]
// @summary Refresh tokens
// @description Refresh tokens (access and refresh token)
// @tags auth
// @param access_token header string true "Access token"
// @param refresh_token header string true "Refresh token"
// @success 200 {object} model.TokenPair "Refreshed Tokens"
// @failure 500 {object} status.Status "Internal server error"
func (v *Handler) Refresh(c *gin.Context) {
	var m model.RefreshRequest
	if err := c.BindHeader(&m); err != nil {
		v.response(c, status.StatusBadRequest)
		return
	}

	ctx, cancel := v.NewContext(c, 5, "device_info")
	defer cancel()

	tokenPair, err := v.service.Auth().Refresh(ctx, m)
	if err != nil {
		v.response(c, v.ParseError(err))
		return
	}

	v.response(c, status.StatusOk.AddData(tokenPair))
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

	ctx, cancel := v.NewContext(c, 5, "device_info")
	defer cancel()

	tokenPair, err := v.service.Auth().Register(ctx, m)
	if err != nil {
		v.response(c, v.ParseError(err))
		return
	}
	v.response(c, status.StatusOk.AddData(tokenPair))
}
