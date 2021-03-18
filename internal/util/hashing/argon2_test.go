package hashing_test

import (
	"regexp"
	"testing"

	"allaboutapps.dev/aw/go-starter/internal/util/hashing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	hashRegex, err := regexp.Compile(`^\$argon2id\$v=19\$m=65536,t=1,p=4\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	require.NoError(t, err, "failed to compile hash regex")

	hash1, err := hashing.HashPassword("t3stp4ssw0rd", hashing.DefaultArgon2Params)
	require.NoError(t, err, "failed to hash password")

	assert.Truef(t, hashRegex.MatchString(hash1), "hash %q is not formatted properly", hash1)

	hash2, err := hashing.HashPassword("t3stp4ssw0rd", hashing.DefaultArgon2Params)
	require.NoError(t, err, "failed to hash password")

	assert.Truef(t, hashRegex.MatchString(hash2), "hash %q is not formatted properly", hash2)

	assert.NotEqualf(t, hash1, hash2, "hashes %q and %q are not unique", hash1, hash2)
}

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := hashing.HashPassword("t3stp4ssw0rd", hashing.DefaultArgon2Params)
		if err != nil {
			b.Errorf("failed to hash password #%d: %v", i, err)
		}
	}
}

func TestComparePasswordAndHash(t *testing.T) {
	hash := "$argon2id$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk"

	match, err := hashing.ComparePasswordAndHash("t3stp4ssw0rd", hash)
	require.NoError(t, err)
	assert.True(t, match)

	match, err = hashing.ComparePasswordAndHash("wr0ngt3stp4ssw0rd", hash)
	require.NoError(t, err)
	assert.False(t, match)
}

func BenchmarkComparePasswordAndHash(b *testing.B) {
	hash := "$argon2id$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk"

	for i := 0; i < b.N; i++ {
		_, err := hashing.ComparePasswordAndHash("t3stp4ssw0rd", hash)
		if err != nil {
			b.Errorf("failed to compare password and hash #%d: %v", i, err)
		}
	}
}

func BenchmarkCompareWrongPasswordAndHash(b *testing.B) {
	hash := "$argon2id$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk"

	for i := 0; i < b.N; i++ {
		_, err := hashing.ComparePasswordAndHash("wr0ngt3stp4ssw0rd", hash)
		if err != nil {
			b.Errorf("failed to compare wrong password and hash #%d: %v", i, err)
		}
	}
}

func TestComparePasswordAndInvalidHash(t *testing.T) {
	_, err := hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2i$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrInvalidArgon2Hash, err)

	_, err = hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrInvalidArgon2Hash, err)

	_, err = hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=20$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrIncompatibleArgon2Version, err)

	_, err = hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=a$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrIncompatibleArgon2Version, err)

	_, err = hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=a,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrInvalidArgon2Hash, err)

	_, err = hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=a,p=4$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrInvalidArgon2Hash, err)

	_, err = hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=1,p=a$c8FqPHMT83tyxE2v0xDAFw$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrInvalidArgon2Hash, err)

	_, err = hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=1,p=4$ä$s2qmbRoRRbfyLIVFUzRwzE7F8PLjchpLKaV7Wf7tHgk")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrInvalidArgon2Hash, err)

	_, err = hashing.ComparePasswordAndHash("t3stp4ssw0rd", "$argon2id$v=19$m=65536,t=1,p=4$c8FqPHMT83tyxE2v0xDAFw$ä")
	require.Error(t, err)
	assert.Equal(t, hashing.ErrInvalidArgon2Hash, err)
}
