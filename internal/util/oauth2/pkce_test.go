package oauth2_test

import (
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util/oauth2"
	"github.com/stretchr/testify/assert"
)

func TestGetPKCECodeChallengeS256(t *testing.T) {
	verifier := "U7MEZRmshzwIHRIGvF5iy6FLKgTtUHV0Vb0Hpczh6jJ_XZKcQIurow2LvsjG6hx2k57s9Pz8UmCZTvazosnniTM-z6EC.skJlQMGA~8ue3LMiOWdFYTfsLdX8GKol285"
	expected := "Jg697bAjhzV1upYvV9R04784OFNVRAZh2IjeFlMJ8bE"

	challenge := oauth2.GetPKCECodeChallengeS256(verifier)
	assert.Equal(t, expected, challenge)
}
