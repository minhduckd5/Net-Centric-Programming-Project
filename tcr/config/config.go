package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the server configuration
type Config struct {
	Server struct {
		Host         string `json:"host"`
		Port         int    `json:"port"`
		ReadTimeout  int    `json:"read_timeout"`
		WriteTimeout int    `json:"write_timeout"`
		IdleTimeout  int    `json:"idle_timeout"`
	} `json:"server"`
	Game struct {
		TickIntervalMs  int    `json:"tick_interval_ms"`
		MatchTimeoutSec int    `json:"match_timeout_sec"`
		MaxPlayers      int    `json:"max_players"`
		LogLevel        string `json:"log_level"`
	} `json:"game"`
	Security struct {
		RateLimit     int    `json:"rate_limit"`
		RateWindowSec int    `json:"rate_window_sec"`
		PasswordSalt  string `json:"password_salt"`
	} `json:"security"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig checks if the configuration values are valid
func validateConfig(config *Config) error {
	// Server validation
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}
	if config.Server.ReadTimeout <= 0 {
		return fmt.Errorf("invalid read timeout: %d", config.Server.ReadTimeout)
	}
	if config.Server.WriteTimeout <= 0 {
		return fmt.Errorf("invalid write timeout: %d", config.Server.WriteTimeout)
	}
	if config.Server.IdleTimeout <= 0 {
		return fmt.Errorf("invalid idle timeout: %d", config.Server.IdleTimeout)
	}

	// Game validation
	if config.Game.TickIntervalMs <= 0 {
		return fmt.Errorf("invalid tick interval: %d", config.Game.TickIntervalMs)
	}
	if config.Game.MatchTimeoutSec <= 0 {
		return fmt.Errorf("invalid match timeout: %d", config.Game.MatchTimeoutSec)
	}
	if config.Game.MaxPlayers <= 0 {
		return fmt.Errorf("invalid max players: %d", config.Game.MaxPlayers)
	}
	if config.Game.LogLevel != "debug" && config.Game.LogLevel != "info" {
		return fmt.Errorf("invalid log level: %s", config.Game.LogLevel)
	}

	// Security validation
	if config.Security.RateLimit <= 0 {
		return fmt.Errorf("invalid rate limit: %d", config.Security.RateLimit)
	}
	if config.Security.RateWindowSec <= 0 {
		return fmt.Errorf("invalid rate window: %d", config.Security.RateWindowSec)
	}
	if config.Security.PasswordSalt == "" {
		return fmt.Errorf("password salt cannot be empty")
	}

	return nil
}
