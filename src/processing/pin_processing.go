package processing

import (
	"fmt"
)

func PinProcessing() (string, string, error) {
	// Generate new pin
	pin, err := PinGenerating()
	if err != nil {
		return "", "", fmt.Errorf("failed to create new pin: %w", err)
	}

	// Hash generated pin
	hash, err := PinHashing(pin)
	if err != nil {
		return "", "", fmt.Errorf("failed to hash pin: %w", err)
	}

	return hash, pin, nil
}
