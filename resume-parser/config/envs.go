package config

import (
	"os"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

type Config struct {
	// TODO: Add configs
	ServiceAccount string
}

func initConfig() Config {
	godotenv.Load()
	return Config{
		ServiceAccount: getEnv("CREDENTIALS", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
