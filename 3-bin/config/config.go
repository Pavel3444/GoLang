package config

import (
	"fmt"
	"os"
)

type Config struct {
	Key string
}

func LoadConfig() (*Config, error) {
	r := os.Getenv("KEY")
	if r == "" {
		return nil, fmt.Errorf("key not found")
	}
	return &Config{Key: r}, nil
}
