package oauth2

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"allaboutapps.dev/aw/go-starter/internal/util"
)

const (
	defaultVerifierLength = 128
)

func GetPKCECodeVerifier() (string, error) {
	// for details regarding possible characters in verifier, see:
	// https://tools.ietf.org/html/rfc7636#section-4.1
	verifier, err := util.GenerateRandomString(defaultVerifierLength, []util.CharRange{util.CharRangeNumeric, util.CharRangeAlphaLowerCase, util.CharRangeAlphaUpperCase}, "-._~")
	if err != nil {
		return "", fmt.Errorf("failed to generate random string: %w", err)
	}

	return verifier, nil
}

func GetPKCECodeChallengeS256(verifier string) string {
	// for details regarding transformation of verifier to challenge see:
	// https://tools.ietf.org/html/rfc7636#section-4.2
	// base64 encoding must be unpadded, URL encoding:
	// https://tools.ietf.org/html/rfc7636#page-17
	sum := sha256.Sum256([]byte(verifier))
	b64 := base64.RawURLEncoding.EncodeToString(sum[:])

	return b64
}
