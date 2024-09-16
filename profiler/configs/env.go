package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

type Config struct {
	// MYSQL Config
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string

	RabbitMQAddr string
	RedisAddr    string
	GithubToken  string
}

func initConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("ERROR LOADING ENVIRONMENT VARIABLES: ", err)
	}
	return Config{
		DBUser:       getEnv("DB_USER", ""),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		DBAddress:    fmt.Sprintf("%s:%s", getEnv("DB_HOST", ""), getEnv("DB_PORT", "")),
		DBName:       getEnv("DB_NAME", ""),
		RabbitMQAddr: getEnv("RABBITMQ_ADDR", ""),
		RedisAddr:    getEnv("REDIS_ADDR", ""),
		GithubToken:  getEnv("GITHUB_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
