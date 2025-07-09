package utils

import (
	"os"
	"github.com/joho/godotenv"
	"fmt"
)

var (
	DB_USERNAME string
	DB_URI string
	DB_PASSWORD string
)

func LoadEnv() (err error) {
	err = godotenv.Load()

	if err != nil {
		fmt.Println(".env cannot be loaded")
		return err
	}

  DB_USERNAME = os.Getenv("DB_USERNAME")

  DB_URI = os.Getenv("DB_URI")

  DB_PASSWORD = os.Getenv("DB_PASSWORD")

	return
}
