package processing

import (
	"math/rand"
	"strconv"
)

func PinGenerating() (string, error) {
	minVal := 100000
	maxVal := 999999
	pin := rand.Intn(maxVal-minVal) + minVal

	return strconv.Itoa(pin), nil
}
