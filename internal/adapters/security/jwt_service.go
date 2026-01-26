package security

import (
	"errors"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secret          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuer          string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// NewJWTService creates a new JWT-based token service
func NewJWTService(config JWTConfig) auth.TokenService {
	return &jwtService{
		secret:          config.Secret,
		accessTokenTTL:  config.AccessTokenTTL,
		refreshTokenTTL: config.RefreshTokenTTL,
		issuer:          config.Issuer,
	}
}

func (j *jwtService) GenerateAccessToken(userID, email, username string) (*auth.GeneratedToken, error) {
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

	return &auth.GeneratedToken{
		Token:     tokenString,
		UserID:    userID,
		Email:     email,
		Username:  username,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
	}, nil
}

func (j *jwtService) GenerateRefreshToken(userID, email, username string) (*auth.GeneratedToken, error) {
	now := time.Now()
	expiresAt := now.Add(j.refreshTokenTTL)

	claims := jwt.MapClaims{
		"user_id":  userID,
		"email":    email,
		"username": username,
		"exp":      expiresAt.Unix(),
		"iss":      j.issuer,
		"iat":      now.Unix(),
		"type":     "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	return &auth.GeneratedToken{
		Token:     tokenString,
		UserID:    userID,
		Email:     email,
		Username:  username,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
	}, nil
}

func (j *jwtService) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
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

	tokenClaims := &auth.TokenClaims{
		UserID:   getStringClaim(claims, "user_id"),
		Email:    getStringClaim(claims, "email"),
		Username: getStringClaim(claims, "username"),
	}

	if iat, ok := claims["iat"].(float64); ok {
		tokenClaims.IssuedAt = time.Unix(int64(iat), 0)
	}
	if exp, ok := claims["exp"].(float64); ok {
		tokenClaims.ExpiresAt = time.Unix(int64(exp), 0)
	}

	return tokenClaims, nil
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
