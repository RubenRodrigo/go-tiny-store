package consts

import "time"

type AuthConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewAuthConfig() AuthConfig {
	return AuthConfig{
		AccessTokenTTL:  24 * time.Hour,
		RefreshTokenTTL: 30 * 24 * time.Hour,
	}
}
