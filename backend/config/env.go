package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

type Config struct {
	PublicHost string
	Port       string

	// MYSQL Config
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string

	AuthSecret    string
	JWTExpiration time.Duration
}

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:    getEnv("PUBLIC_HOST", "http://localhost"),
		Port:          getEnv("PORT", ":8000"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", "pwd"),
		DBAddress:     fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", ":3306")),
		DBName:        getEnv("DB_NAME", "test"),
		AuthSecret:    getEnv("AUTH_SECRET", "random_string"),
		JWTExpiration: time.Duration(time.Hour * 24 * 30), // 30 Days
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
