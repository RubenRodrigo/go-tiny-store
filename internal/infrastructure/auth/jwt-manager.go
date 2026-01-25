package auth

import (
	"errors"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

func NewJWTManager(config JWTConfig) TokenService {
	return &JWTManager{
		secret:          config.Secret,
		accessTokenTTL:  config.AccessTokenTTL,
		refreshTokenTTL: config.RefreshTokenTTL,
		issuer:          config.Issuer,
	}
}

func (j *JWTManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrAuthTokenInvalid
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, apperrors.ErrAuthTokenInvalid
	}

	if !token.Valid {
		return nil, apperrors.ErrAuthTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	// Parse claims into our struct
	tokenClaims := &TokenClaims{
		UserID:   getStringClaim(claims, "user_id"),
		Email:    getStringClaim(claims, "email"),
		Username: getStringClaim(claims, "username"),
	}

	// Parse timestamps
	if iat, ok := claims["iat"].(float64); ok {
		tokenClaims.IssuedAt = time.Unix(int64(iat), 0)
	}
	if exp, ok := claims["exp"].(float64); ok {
		tokenClaims.ExpiresAt = time.Unix(int64(exp), 0)
	}

	return tokenClaims, nil
}

func (j *JWTManager) GenerateAccessToken(userID, email, username string) (*TokenGeneratedClaims, error) {
	now := time.Now()
	expiresAt := now.Add(j.accessTokenTTL)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"username": username,
		"exp":      expiresAt.Unix(),
		"iss":      j.issuer,
		"iat":      now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	return &TokenGeneratedClaims{
		Token:     tokenString,
		Username:  username,
		UserID:    userID,
		Email:     email,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
	}, nil
}

// GenerateRefreshToken creates a JWT-based refresh token with user information
func (j *JWTManager) GenerateRefreshToken(userID, email, username string) (*TokenGeneratedClaims, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshTokenTTL)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"username": username,
		"exp":      expiresAt.Unix(),
		"iss":      j.issuer,
		"iat":      now.Unix(),
		"type":     "refresh", // Distinguish from access tokens
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	return &TokenGeneratedClaims{
		Token:     tokenString,
		Username:  username,
		UserID:    userID,
		Email:     email,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
	}, nil
}

// Helper function to safely extract string claims
func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
