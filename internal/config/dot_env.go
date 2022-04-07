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
