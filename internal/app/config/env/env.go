package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Get(name string) string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	val := os.Getenv(name)
	return val
}
