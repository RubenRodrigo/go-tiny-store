package lib

import (
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
}

func NewJWTManager(secret []byte) *JWTManager {
	return &JWTManager{
		secret: secret,
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

func (j *JWTManager) GenerateToken(userID uint, email, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}
