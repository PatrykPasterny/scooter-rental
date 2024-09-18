package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	HTTP  int    `env:"HTTP,required"`
	Name  string `env:"NAME,required"`
	Users string `env:"USERS,required"`
}

func NewConfig(ctx context.Context, configPath string) (*Config, error) {
	if err := godotenv.Load(configPath); err != nil {
		return nil, fmt.Errorf("loading config files: %w", err)
	}

	var c Config

	if err := envconfig.Process(ctx, &c); err != nil {
		return nil, fmt.Errorf("processing environment config: %w", err)
	}

	return &c, nil
}

func (c *Config) GetUsersMap() map[string]bool {
	result := make(map[string]bool)

	userIDs := strings.Split(c.Users, ",")

	for _, userID := range userIDs {
		result[userID] = true
	}

	return result
}
