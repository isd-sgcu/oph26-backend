package initializer

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	if os.Getenv("DOCKER_CONTAINER") == "true" || os.Getenv("COMPOSE_PROJECT_NAME") != "" {
		log.Println("Running in containerized environment, using environment variables from docker-compose")
		return
	}

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println(".env file loaded successfully")
	}
}