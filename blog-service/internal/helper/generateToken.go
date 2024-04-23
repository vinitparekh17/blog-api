package helper

import (
	"crypto/rand"
	"fmt"
)

func GenerateToken() (token string, err error) {
	randBytes := make([]byte, 16)
	_, err = rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", randBytes), nil
}
