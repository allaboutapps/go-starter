package pgconsumer

import "os"

type ClientConfig struct {
	BaseURL    string
	APIVersion string
}

func DefaultClientConfigFromEnv() ClientConfig {
	return ClientConfig{
		BaseURL: getEnv("PGCONSUMER_BASE_URL", "http://127.0.0.1:8080/api"),
		// BaseURL:    getEnv("PGCONSUMER_BASE_URL", "http://pgserve:8080/api"),
		APIVersion: getEnv("PGCONSUMER_API_VERSION", "v1"),
	}
}

func getEnv(key string, fallback string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return v
}
