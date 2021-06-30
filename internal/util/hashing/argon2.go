package hashing

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Inspired by: https://github.com/alexedwards/argon2id @ 2020-04-22T14:13:23ZZ

const (
	// Argon2HashID represents the hash ID set in the (pseudo) modular crypt format used to store the hashed password and params in a single string.
	Argon2HashID = "argon2id"
)

var (
	// ErrInvalidArgon2Hash indicates the argon2id hash was malformed and could not be decoded.
	ErrInvalidArgon2Hash = errors.New("invalid argon2id hash")
	// ErrIncompatibleArgon2Version indicates the argon2id hash provided was generated with a different, incompatible argon2 version.
	ErrIncompatibleArgon2Version = errors.New("incompatible argon2 version")
)

func HashPassword(password string, params *Argon2Params) (hash string, err error) {
	salt, err := generateSalt(params.SaltLength)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s", Argon2HashID, argon2.Version, params.Memory, params.Time, params.Threads, b64Salt, b64Key), nil
}

func ComparePasswordAndHash(password string, hash string) (matches bool, err error) {
	params, salt, key, err := decodeArgon2Hash(hash)
	if err != nil {
		return false, err
	}

	pKey := argon2.IDKey([]byte(password), salt, params.Time, params.Memory, params.Threads, params.KeyLength)

	if subtle.ConstantTimeEq(int32(len(key)), int32(len(pKey))) == 0 {
		return false, nil
	}

	if subtle.ConstantTimeCompare(key, pKey) == 0 {
		return false, nil
	}

	return true, nil
}

func decodeArgon2Hash(hash string) (params *Argon2Params, salt []byte, key []byte, err error) {
	vals := strings.Split(hash, "$") // splits into array of 6 values, with val[0] being empty --> length/indicies "offset" by one
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidArgon2Hash
	}
	if vals[1] != Argon2HashID {
		return nil, nil, nil, ErrInvalidArgon2Hash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, ErrIncompatibleArgon2Version
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleArgon2Version
	}

	params = &Argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Time, &params.Threads)
	if err != nil {
		return nil, nil, nil, ErrInvalidArgon2Hash
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, ErrInvalidArgon2Hash
	}

	key, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, ErrInvalidArgon2Hash
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
