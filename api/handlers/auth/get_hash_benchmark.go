package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"allaboutapps.at/aw/go-mranftl-sample/api"
	"allaboutapps.at/aw/go-mranftl-sample/pkg/auth/hashing"
	"github.com/labstack/echo/v4"
)

func GetHashBenchmarkRoute(s *api.Server) *echo.Route {
	return s.Router.ApiV1Auth.GET("/hash/benchmark", getHashBenchmarkHandler(s))
}

func getHashBenchmarkHandler(s *api.Server) echo.HandlerFunc {

	return func(c echo.Context) error {
		count, err := strconv.Atoi(c.QueryParam("count"))
		if err != nil {
			count = 1
		}

		params := hashing.DefaultArgon2ParamsFromEnv()

		memory, err := strconv.ParseUint(c.QueryParam("memory"), 10, 32)
		if err == nil {
			params.Memory = uint32(memory)
		}
		threads, err := strconv.ParseUint(c.QueryParam("threads"), 10, 8)
		if err == nil {
			params.Threads = uint8(threads)
		}
		t, err := strconv.ParseUint(c.QueryParam("time"), 10, 32)
		if err == nil {
			params.Time = uint32(t)
		}

		totalStart := time.Now()
		for i := 0; i < count; i++ {
			start := time.Now()
			_, err := hashing.HashPassword("t3stp4ssw0rd", params)
			if err != nil {
				return err
			}
			fmt.Printf("hash #%d: %s\n", i, time.Since(start))
		}
		fmt.Printf("total #%d: %s\n", count, time.Since(totalStart))

		return c.NoContent(http.StatusNoContent)
	}
}
