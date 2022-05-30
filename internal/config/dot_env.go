package config

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/subosito/gotenv"
)

// os.SetEnv func signature
type envSetter = func(key string, value string) error

// DotEnvTryLoad forcefully overrides ENV variables through **a maybe available** .env file.
//
// This function will always remain silent if a .env file does not exist!
// If we successfully apply an ENV file, we will log a warning.
// If there are any other errors, we will panic!
//
// This mechanism should only be used **locally** to easily inject (gitignored)
// secrets into your ENV. Non-existing .env files are actually the **best case**.
//
// When running normally (not within tests):
// DotEnvTryLoad("/path/tp/my.env.local", os.SetEnv)
//
// For tests (and autoreset) use t.Setenv:
// DotEnvTryLoad("/path/to/my.env.test.local", func(k string, v string) error { t.Setenv(k, v); return nil })
func DotEnvTryLoad(absolutePathToEnvFile string, setEnvFn envSetter) {
	err := DotEnvLoad(absolutePathToEnvFile, setEnvFn)

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic().Err(err).Str("envFile", absolutePathToEnvFile).Msg(".env parse error!")
		}
	} else {
		log.Warn().Str("envFile", absolutePathToEnvFile).Msg(".env overrides ENV variables!")
	}
}

// DotEnvLoad forcefully overrides ENV variables through the supplied .env file.
//
// This mechanism should only be used **locally** to easily inject (gitignored)
// secrets into your ENV.
//
// When running normally (not within tests):
// DotEnvLoad("/path/to/my.env.local", os.SetEnv)
//
// For tests (and ENV var autoreset) use t.Setenv:
// DotEnvLoad("/path/to/my.env.test.local", func(k string, v string) error { t.Setenv(k, v); return nil })
func DotEnvLoad(absolutePathToEnvFile string, setEnvFn envSetter) error {

	file, err := os.Open(absolutePathToEnvFile)

	if err != nil {
		return err
	}

	defer file.Close()

	envs, err := gotenv.StrictParse(file)

	if err != nil {
		return err
	}

	for key, value := range envs {
		if err := setEnvFn(key, value); err != nil {
			return err
		}
	}

	return nil
}
