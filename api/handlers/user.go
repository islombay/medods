package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"medods/api/status"
	"time"
)

// GetMe
// @id getme
// @router	/api/me [get]
// @summary Get Information about current user
// @security ApiKeyAuth
// @tags user
// @success 200 {object} model.User "User"
// @failure 500 {object} status.Status "Internal server error"
func (v *Handler) GetMe(c *gin.Context) {
	user_id, ok := GetValue[string](c, "user_id")
	if !ok {
		v.response(c, status.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user, err := v.service.User().GetUser(ctx, user_id)
	if err != nil {
		v.response(c, v.ParseError(err))
		return
	}

	v.response(c, status.StatusOk.AddData(user))
}
