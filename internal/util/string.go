package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

// GenerateRandomBytes returns n random bytes securely generated using the system's default CSPRNG.
//
// An error will be returned if reading from the secure random number generator fails, at which point
// the returned result should be discarded and not used any further.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomBase64String returns a string with n random bytes securely generated using the system's
// default CSPRNG in base64 encoding. The resulting string might not be of length n as the encoding for
// the raw bytes generated may vary.
//
// An error will be returned if reading from the secure random number generator fails, at which point
// the returned result should be discarded and not used any further.
func GenerateRandomBase64String(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

// GenerateRandomHexString returns a string with n random bytes securely generated using the system's
// default CSPRNG in hexadecimal encoding. The resulting string might not be of length n as the encoding
// for the raw bytes generated may vary.
//
// An error will be returned if reading from the secure random number generator fails, at which point
// the returned result should be discarded and not used any further.
func GenerateRandomHexString(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
