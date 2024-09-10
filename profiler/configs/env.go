package configs

import (
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
}

func initConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("ERROR LOADING ENVIRONMENT VARIABLES: ", err)
	}
	return Config{
		RabbitMQAddr: getEnv("RABBITMQ_ADDR", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
