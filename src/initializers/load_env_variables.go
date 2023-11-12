package initializers

import (
	"fmt"
	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	err := godotenv.Load()

	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	return nil
}
