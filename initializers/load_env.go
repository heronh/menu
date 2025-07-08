package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file, relying on system environment variables.")
	} else {
		log.Println(".env file loaded successfully.")
	}
}
