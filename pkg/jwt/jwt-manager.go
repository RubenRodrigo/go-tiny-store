package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/pkg/consts"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	cfg    consts.AuthConfig
}

func NewJWTManager(secret []byte, cfg consts.AuthConfig) *JWTManager {
	return &JWTManager{
		secret: secret,
		cfg:    cfg,
	}
}

func (j *JWTManager) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrAuthTokenInvalid
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, apperrors.ErrAuthTokenInvalid
	}

	if !token.Valid {
		return nil, apperrors.ErrAuthTokenInvalid
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, apperrors.ErrAuthTokenExpired
		}
	}

	return claims, nil
}

func (j *JWTManager) GenerateAccessToken(userID, email, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(j.cfg.AccessTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

// GenerateRefreshToken creates a random, URL-safe token string
func (m *JWTManager) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	rand.Read(bytes)

	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}
