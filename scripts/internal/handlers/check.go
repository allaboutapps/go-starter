//go:build scripts

// This program checks /internal/api/handlers.go and /internal/types/spec_handlers.go
// It can be invoked by running go run -tags scripts scripts/handlers/check_handlers.go

// Supported args:
// --print-all

package handlers

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"allaboutapps.dev/aw/go-starter/internal/api"
	"allaboutapps.dev/aw/go-starter/internal/api/router"
	"allaboutapps.dev/aw/go-starter/internal/config"
	"allaboutapps.dev/aw/go-starter/internal/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func CheckHandlers(printAll bool) error {
	// we initialize a minimal echo server without any other deps
	// so we can attach the current defined routes and read them
	log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "15:04:05"
	}))

	defaultConfig := config.DefaultServiceConfigFromEnv()
	defaultConfig.Echo.ListenAddress = ":0"

	s := api.NewServer(defaultConfig)
	err := router.Init(s)
	if err != nil {
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	// swaggerspec vs routes
	routes := s.Router.Routes
	swaggerSpec := types.NewSwaggerSpec()

	var wg sync.WaitGroup

	for _, route := range routes {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// replace named echo ":param" to swagger path params "{param}" (curly braces) to properly match paths
			fragments := strings.Split(route.Path, "/")

			for i, fragment := range fragments {
				if strings.HasPrefix(fragment, ":") {
					fragments[i] = "{" + strings.TrimLeft(fragment, ":") + "}"
				}
			}

			swaggerPath := strings.Join(fragments, "/")

			ok := swaggerSpec.Handlers[route.Method][swaggerPath]

			if !ok {
				log.Warn().Msgf("%s %s\n    WARNING: Missing swagger spec in api/swagger.yml!", route.Method, route.Path)
			} else if printAll {
				log.Info().Msgf("%s %s", route.Method, route.Path)
			}
		}()
	}

	for method, v := range swaggerSpec.Handlers {
		for path := range v {
			// NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
			ttMethod := method
			ttPath := path

			wg.Add(1)

			go func() {
				defer wg.Done()

				// replace named swagger path params "{param}" to echo ":param" (dotted) to properly match paths
				fragments := strings.Split(ttPath, "/")

				for i, fragment := range fragments {
					if strings.HasPrefix(fragment, "{") && strings.HasSuffix(fragment, "}") {
						fragments[i] = ":" + strings.TrimLeft(strings.TrimRight(fragment, "}"), "{")
					}
				}

				echoPath := strings.Join(fragments, "/")

				hasMatch := false

				for _, route := range routes {
					if route.Method == ttMethod && route.Path == echoPath {
						hasMatch = true
						break
					}
				}

				if !hasMatch {
					log.Warn().Msgf("%s %s\n    WARNING: Missing route implementation in internal/api/handlers/*!", ttMethod, echoPath)
				}
			}()
		}
	}

	// we have our routes, the server is no longer needed.
	if errs := s.Shutdown(context.Background()); len(errs) > 0 {
		log.Error().Errs("shutdownErrors", errs).Msg("Failed to stop introspection server")
	}

	wg.Wait()

	return nil
}
