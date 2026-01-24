package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ExtractStringClaim(claims jwt.MapClaims, key string) (string, error) {
	val, ok := claims[key]
	if !ok {
		return "", fmt.Errorf("claim %s not found", key)

	}

	strVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("claim %s is not a string", key)
	}

	return strVal, nil
}
