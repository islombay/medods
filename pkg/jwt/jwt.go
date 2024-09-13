package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"medods/internal/error_list"
	"os"
	"strconv"
	"time"
)

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserID    string `json:"user_id"`
	IP        string `json:"ip"`
	SessionID string `json:"session_id"`
}

func Generate(userID, ip, session_id string) (string, error) {
	// Will convert
	expire_duration_minutes, _ := strconv.Atoi(os.Getenv("TOKEN_ACCESS_DURATION_MINUTES"))

	expiresAt := time.Now().Add(time.Duration(expire_duration_minutes) * time.Minute)

	claims := &AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserID:    userID,
		IP:        ip,
		SessionID: session_id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString([]byte(os.Getenv("TOKEN_SECRET_KEY")))
}

func GenerateRefreshToken() (string, error) {
	// randomize and make unique refresh token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	refreshToken := base64.StdEncoding.EncodeToString(tokenBytes)

	return refreshToken, nil
}

func ParseToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET_KEY")), nil
	})

	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return nil, err
		}
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, error_list.TokenExpired
		}
		return nil, err
	}

	return claims, nil
}

func GeneratePairTokens(userID, ip, session_id string) (string, string, error) {
	accessToken, err := Generate(userID, ip, session_id)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
