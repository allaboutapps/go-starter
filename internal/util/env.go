package util

import (
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

const (
	mgmtSecretLen = 16
)

var (
	mgmtSecret     string
	mgmtSecretOnce sync.Once
)

func GetEnv(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func GetEnvEnum(key string, defaultVal string, allowedValues []string) string {
	if !ContainsString(allowedValues, defaultVal) {
		log.Panic().Str("key", key).Str("value", defaultVal).Msg("Default value is not in the allowed values list.")
	}

	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}

	if !ContainsString(allowedValues, val) {
		log.Error().Str("key", key).Str("value", val).Msg("Value is not allowed. Fallback to default value.")
		return defaultVal
	}

	return val
}

func GetEnvAsInt(key string, defaultVal int) int {
	strVal := GetEnv(key, "")

	if val, err := strconv.Atoi(strVal); err == nil {
		return val
	}

	return defaultVal
}

func GetEnvAsUint32(key string, defaultVal uint32) uint32 {
	strVal := GetEnv(key, "")

	if val, err := strconv.ParseUint(strVal, 10, 32); err == nil {
		return uint32(val)
	}

	return defaultVal
}

func GetEnvAsUint8(key string, defaultVal uint8) uint8 {
	strVal := GetEnv(key, "")

	if val, err := strconv.ParseUint(strVal, 10, 8); err == nil {
		return uint8(val)
	}

	return defaultVal
}

func GetEnvAsBool(key string, defaultVal bool) bool {
	strVal := GetEnv(key, "")

	if val, err := strconv.ParseBool(strVal); err == nil {
		return val
	}

	return defaultVal
}

func GetEnvAsStringArr(key string, defaultVal []string, separator ...string) []string {
	strVal := GetEnv(key, "")

	if len(strVal) == 0 {
		return defaultVal
	}

	sep := ","
	if len(separator) >= 1 {
		sep = separator[0]
	}

	return strings.Split(strVal, sep)
}

func GetEnvAsURL(key string, defaultVal string) *url.URL {
	strVal := GetEnv(key, "")

	if len(strVal) == 0 {
		u, err := url.Parse(defaultVal)
		if err != nil {
			log.Panic().Str("key", key).Str("defaultVal", defaultVal).Err(err).Msg("Failed to parse default value for env variable as URL")
		}

		return u
	}

	u, err := url.Parse(strVal)
	if err != nil {
		log.Panic().Str("key", key).Str("strVal", strVal).Err(err).Msg("Failed to parse env variable as URL")
	}

	return u
}

// GetMgmtSecret returns the management secret for the app server, mainly used by health check and readiness endpoints.
// It first attempts to retrieve a value from the given environment variable and generates a cryptographically secure random string
// should no env var have been set.
// Failure to generate a random string will cause a panic as secret security cannot be guaranteed otherwise.
// Subsequent calls to GetMgmtSecret during the server's runtime will always return the same randomly generated secret for consistency.
func GetMgmtSecret(envKey string) string {
	val := GetEnv(envKey, "")

	if len(val) > 0 {
		return val
	}

	mgmtSecretOnce.Do(func() {
		var err error
		mgmtSecret, err = GenerateRandomHexString(mgmtSecretLen)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to generate random management secret")
		}

		log.Warn().Str("envKey", envKey).Str("mgmtSecret", mgmtSecret).Msg("Could not retrieve management secret from env key, using randomly generated one")
	})

	return mgmtSecret
}
