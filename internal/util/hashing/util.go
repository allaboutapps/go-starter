package hashing

import (
	"crypto/rand"
	"fmt"
)

func generateSalt(n uint32) ([]byte, error) {
	result := make([]byte, n)

	_, err := rand.Read(result)
	if err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	return result, nil
}
