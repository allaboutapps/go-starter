package router

import "github.com/rs/zerolog"

type echoLogger struct {
	level zerolog.Level
	log   zerolog.Logger
}

func (l *echoLogger) Write(p []byte) (int, error) {
	l.log.WithLevel(l.level).Msgf("%s", p)
	return len(p), nil
}
