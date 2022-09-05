package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Get(name string) string {

	val := os.Getenv(name)
	if val == "" {
		return getFromFile(name)
	}
	return val
}

func getFromFile(name string) string {
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}

	val := os.Getenv(name)
	return val
}
