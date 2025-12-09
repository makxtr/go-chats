package utils

import (
	"errors"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func GetByteFromStr(str string) (byte, error) {
	if str == "" {
		return 0, errors.New("input string is empty")
	}

	numInt, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("failed to parse string '%s' to integer: %w", str, err)
	}

	if numInt < 0 || numInt > 255 {
		return 0, fmt.Errorf("value %d from string '%s' is out of byte range (0-255)", numInt, str)
	}

	return byte(numInt), nil
}

func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}
