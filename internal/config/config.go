package config

import "os"

// Config holds runtime configuration values.
type Config struct {
	Port        string
	ComputeAddr string
}

// Load reads configuration from environment variables.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	computeAddr := os.Getenv("COMPUTE_ADDR")
	if computeAddr == "" {
		computeAddr = "localhost:50051"
	}

	return Config{
		Port:        port,
		ComputeAddr: computeAddr,
	}
}
