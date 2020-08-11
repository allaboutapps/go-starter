package hashing

import "allaboutapps.dev/aw/go-starter/internal/util"

var (
	// DefaultArgon2Params represents Argon2ID parameter recommendations in accordance with:
	// https://pkg.go.dev/golang.org/x/crypto@v0.0.0-20200420201142-3c4aac89819a/argon2?tab=doc#IDKey @ 2020-04-22T11:23:38Z
	DefaultArgon2Params = &Argon2Params{
		Time:       1,         // 1 second
		Memory:     64 * 1024, // ~64MB memory costs
		Threads:    4,         // 4 threads
		KeyLength:  32,        // 256 bit key length
		SaltLength: 16,        // 126 bit salt length
	}
)

type Argon2Params struct {
	Time       uint32
	Memory     uint32
	Threads    uint8
	KeyLength  uint32
	SaltLength uint32
}

func DefaultArgon2ParamsFromEnv() *Argon2Params {
	p := &Argon2Params{
		Time:       util.GetEnvAsUint32("AUTH_HASHING_ARGON2_TIME", DefaultArgon2Params.Time),
		Memory:     util.GetEnvAsUint32("AUTH_HASHING_ARGON2_MEMORY", DefaultArgon2Params.Memory),
		Threads:    util.GetEnvAsUint8("AUTH_HASHING_ARGON2_THREADS", DefaultArgon2Params.Threads),
		KeyLength:  util.GetEnvAsUint32("AUTH_HASHING_ARGON2_KEY_LENGTH", DefaultArgon2Params.KeyLength),
		SaltLength: util.GetEnvAsUint32("AUTH_HASHING_ARGON2_SALT_LENGTH", DefaultArgon2Params.SaltLength),
	}

	return p
}
