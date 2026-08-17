package config

import (
	"go-income-expense-tracker-app/internal/constant"
	"go-income-expense-tracker-app/internal/middleware"
	"go-income-expense-tracker-app/internal/utils"
	"strconv"
)

const defaultJWTSecretKey = "dev-secret-key-change-me"

func LoadJWTConfig() middleware.JWTConfig {
	secretKey := utils.GetConfig(constant.JWT_SECRET_KEY)
	if secretKey == "" {
		secretKey = defaultJWTSecretKey
	}

	expiresDuration := 60
	if value := utils.GetConfig(constant.JWT_EXPIRE_DURATION); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			expiresDuration = parsed
		}
	}

	return middleware.JWTConfig{
		SecretKey:       secretKey,
		ExpiresDuration: expiresDuration,
	}
}
