package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/test-go/testify/assert"
	"os"
	"testing"
)

type TestParseTokenCase struct {
	tokenString    string
	needErr        bool
	expectedClaims *AccessTokenClaims
	key            string
}

func TestGenerate(t *testing.T) {
	cases := []struct {
		env        map[string]string
		userID     string
		ip         string
		session_id string
	}{
		{
			env: map[string]string{
				"TOKEN_ACCESS_DURATION_MINUTES": "1",
				"TOKEN_SECRET_KEY":              "secret",
			},
			userID:     "2352",
			ip:         "23.15.3.246",
			session_id: "1234tfdwji",
		},
		{
			env: map[string]string{
				"TOKEN_ACCESS_DURATION_MINUTES": "7",
				"TOKEN_SECRET_KEY":              "",
			},
			userID:     "2352",
			ip:         "23.15.3.246",
			session_id: "",
		},
	}

	for _, c := range cases {
		for k, v := range c.env {
			os.Setenv(k, v)
		}

		tokenString, err := Generate(c.userID, c.ip, c.session_id)
		assert.NoError(t, err, fmt.Sprintf("received error %s", err))

		token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(c.env["TOKEN_SECRET_KEY"]), nil
		})

		assert.NoError(t, err, fmt.Sprintf("received error %s", err))
		assert.True(t, token.Valid, "must be true")

		claims, ok := token.Claims.(*AccessTokenClaims)
		assert.True(t, ok, "Must get claims")

		assert.Equal(t, claims.UserID, c.userID)
		assert.Equal(t, claims.SessionID, c.session_id)
		assert.Equal(t, claims.IP, c.ip)
	}
}
