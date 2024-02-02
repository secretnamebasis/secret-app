package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config struct to hold configuration parameters
type Server struct {
	Port int
}

// Config func to get env value from key
func Config(key string) string {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return ""
	}

	return os.Getenv(key)
}
