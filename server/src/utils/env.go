package utils

import (
	"os"
	"github.com/joho/godotenv"
	// "fmt"
)

var (
	REDIS_USERNAME string
	DB_URI string
	REDIS_PASSWORD string
	REDIS_URI string
	PORT string
	JWT_SECRET, REFRESH_SECRET []byte
)

func LoadEnv() (err error) {
	err = godotenv.Load()

  REDIS_USERNAME = os.Getenv("REDIS_USERNAME")

  DB_URI = os.Getenv("DB_URI")

	if port := os.Getenv("PORT"); port != "" {
		PORT = port
	} else {
		PORT = ":3030"
	}

	REDIS_URI = os.Getenv("REDIS_URI")

  REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")

	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

	REFRESH_SECRET = []byte(os.Getenv("REFRESH_SECRET"))

	return
}
