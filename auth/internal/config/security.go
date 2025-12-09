package config

import (
	"auth/internal/utils"
	"fmt"
)

const (
	refreshTokenSecretKey  = "REFRESH_TOKEN_SECRET_KEY"
	accessTokenSecretKey   = "ACCESS_TOKEN_SECRET_KEY"
	refreshTokenExpiration = "REFRESH_TOKEN_EXP"
	accessTokenExpiration  = "ACCESS_TOKEN_EXP"
)

type SecurityConfig interface {
	RefreshKey() string
	AccessKey() string
	RefreshExp() byte
	AccessExp() byte
}

type securityConfig struct {
	refreshKey string
	accessKey  string
	refreshExp byte
	accessExp  byte
}

func NewSecurityConfig() (SecurityConfig, error) {
	refreshKey, err := GetRequiredEnv(refreshTokenSecretKey)
	if err != nil {
		return nil, err
	}
	accessKey, err := GetRequiredEnv(accessTokenSecretKey)
	if err != nil {
		return nil, err
	}

	refreshExpCfg, err := GetRequiredEnv(refreshTokenExpiration)
	if err != nil {
		return nil, err
	}
	refreshExp, err := utils.GetByteFromStr(refreshExpCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration for %s: %w", refreshTokenExpiration, err)
	}

	accessExpCfg, err := GetRequiredEnv(accessTokenExpiration)
	if err != nil {
		return nil, err
	}
	accessExp, err := utils.GetByteFromStr(accessExpCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expiration for %s: %w", accessTokenExpiration, err)
	}

	return &securityConfig{
		refreshKey: refreshKey,
		accessKey:  accessKey,
		refreshExp: refreshExp,
		accessExp:  accessExp,
	}, nil
}

func (cfg *securityConfig) RefreshKey() string {
	return cfg.refreshKey
}

func (cfg *securityConfig) AccessKey() string {
	return cfg.accessKey
}

func (cfg *securityConfig) RefreshExp() byte {
	return cfg.refreshExp
}

func (cfg *securityConfig) AccessExp() byte {
	return cfg.accessExp
}
