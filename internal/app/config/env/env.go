package env

import (
	"os"

	"github.com/joho/godotenv"
)

func Get(name string) string {
	err := godotenv.Load()
	if err != nil {
		// log.Fatal("Error loading .env file")
	}

	val := os.Getenv(name)
	return val
}
