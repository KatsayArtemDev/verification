package processing

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func PinHashing(pin string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash pin: %v", err)
	}

	return string(hash), nil
}
