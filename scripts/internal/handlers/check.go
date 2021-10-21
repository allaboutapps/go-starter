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
)

func CheckHandlers(printAll bool) {
	// https://golangbyexample.com/print-output-text-color-console/
	// https://gist.github.com/ik5/d8ecde700972d4378d87
	warningLine := "\033[1;33m%s\033[0m\n"

	// we initialize a minimal echo server without any other deps
	// so we can attach the current defined routes and read them
	zerolog.SetGlobalLevel(zerolog.Disabled)
	defaultConfig := config.DefaultServiceConfigFromEnv()
	defaultConfig.Echo.ListenAddress = ":0"

	s := api.NewServer(defaultConfig)
	router.Init(s)

	// swaggerspec vs routes
	routes := s.Router.Routes
	swaggerSpec := types.NewSwaggerSpec()

	var wg sync.WaitGroup

	for _, tt := range routes {
		tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		wg.Add(1)
		go func() {

			defer wg.Done()

			// replace named echo ":param" to swagger path params "{param}" (curly braces) to properly match paths
			fragments := strings.Split(tt.Path, "/")

			for i, fragment := range fragments {
				if strings.HasPrefix(fragment, ":") {
					fragments[i] = "{" + strings.TrimLeft(fragment, ":") + "}"
				}
			}

			swaggerPath := strings.Join(fragments, "/")

			ok := swaggerSpec.Handlers[tt.Method][swaggerPath]

			if !ok {
				fmt.Printf(warningLine, tt.Method+" "+tt.Path+"\n    WARNING: Missing swagger spec in api/swagger.yml!")
			} else {
				if printAll {
					fmt.Println(tt.Method + " " + tt.Path)
				}
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
					fmt.Printf(warningLine, ttMethod+" "+echoPath+"\n    WARNING: Missing route implementation in internal/api/handlers/*!")
				}

			}()
		}
	}

	// we have our routes, the server is no longer needed.
	if err := s.Shutdown(context.Background()); err != nil {
		fmt.Println("Failed to stop introspection server")
	}

	wg.Wait()
}
