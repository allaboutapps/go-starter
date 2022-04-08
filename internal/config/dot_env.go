package config

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/subosito/gotenv"
)

// overrideEnv forcefully overrides ENV variables through the supplied .env file.
//
// This mechanism should only be used **locally** to easily inject (gitignored)
// secrets into your ENV. Non-existing .env files are actually the best case.
// Thus, this function will always remain silent if a .env file does not exist!
//
// If we successfully apply an ENV file, we will log a warning.
// If there are any other errors, we will panic!
func overrideEnv(absolutePathToEnvFile string) {
	err := gotenv.OverLoad(absolutePathToEnvFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Panic().Err(err).Str("envFile", absolutePathToEnvFile).Msg(".env parse error!")
		}
	} else {
		log.Warn().Str("envFile", absolutePathToEnvFile).Msg(".env overrides ENV variables!")
	}
}

// SetEnvFromFile loads a dotenv file and executes the passed setEnvFn per ENV key/value.
// When running in test and you simply want to override the current ENV, simply pass:
//
// SetEnvFromFile("/path/tp/my.test.env.local", t.SetEnv)
//
// This ensures the ENV vars will reset to their original state after your test is finished.
func SetEnvFromFile(absolutePathToEnvFile string, setEnvFn func(key string, val string)) error {

	file, err := os.Open(absolutePathToEnvFile)

	if err != nil {
		return err
	}

	defer file.Close()

	envs, err := gotenv.StrictParse(file)

	if err != nil {
		return err
	}

	for key, val := range envs {
		// if err := setEnvFn(key, val); err != nil {
		// 	return err
		// }
		setEnvFn(key, val)
	}

	return nil
}
