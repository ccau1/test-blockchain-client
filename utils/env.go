package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// get env name
	env := os.Getenv("ENV")
	if env == "" {
		// set default env = env
		env = "dev"
	}
	// load .env.{env}
	godotenv.Load(".env." + env)
	// load original .env
	godotenv.Load()
}
