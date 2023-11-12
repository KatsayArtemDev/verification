package processing

import (
	"golang.org/x/crypto/bcrypt"
)

func PinComparing(correctHash, pin string) error {
	err := bcrypt.CompareHashAndPassword([]byte(correctHash), []byte(pin))
	if err != nil {
		return err
	}
	return nil
}
