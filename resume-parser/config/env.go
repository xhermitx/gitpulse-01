package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Envs = initConfig()

type Config struct {
	PublicHost string
	Port       string

	// TODO: Add configs
	ServiceAccount string

	// For Testing
	TestFolderId string
}

func initConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		return Config{}
	}
	return Config{
		PublicHost:     getEnv("PUBLIC_HOST", ""),
		Port:           getEnv("PORT", ""),
		ServiceAccount: getEnv("CREDENTIALS_JSON", ""),
		TestFolderId:   getEnv("FOLDER_ID", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return fallback
}
