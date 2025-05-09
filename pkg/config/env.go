package config

import (
	"os"
	"strings"
)

const (
	EnvDevelopment = "dev"
	EnvProduction  = "prod"
)

func GetEnvironment() string {
	env := os.Getenv("ENV")

	switch strings.ToLower(env) {
	case "development", "dev":
		return EnvDevelopment
	case "production", "prod":
		return EnvProduction
	default:
		return EnvDevelopment
	}
}

func isDevelopment() bool {
	return GetEnvironment() == EnvDevelopment
}

func isProduction() bool {
	return GetEnvironment() == EnvProduction
}
