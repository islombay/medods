package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func IsUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ParseIP(c *gin.Context) string {
	if c.GetHeader("X-Real-IP") != "" {
		return c.GetHeader("X-Real-IP")
	} else if c.GetHeader("X-Forwarded-For") != "" {
		return c.GetHeader("X-Forwarded-For")
	}
	return c.ClientIP()
}
