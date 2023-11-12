package processing

import (
	"fmt"
	"net/mail"
)

func EmailValidation(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("error when parsing email: %w", err)
	}

	return nil
}
