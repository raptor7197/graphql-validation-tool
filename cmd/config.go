package cmd

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Type     string `yaml:"type"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		DBName   string `yaml:"dbname"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Production bool `yaml:"production"`
}

// LoadConfig reads and parses the config file, with environment variable overrides
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	// Override with environment variables if set
	config.Database.Host = getEnv("DB_HOST", config.Database.Host)
	config.Database.DBName = getEnv("DB_NAME", config.Database.DBName)
	config.Database.User = getEnv("DB_USER", config.Database.User)
	config.Database.Password = getEnv("DB_PASSWORD", config.Database.Password)
	config.Database.SSLMode = getEnv("DB_SSLMODE", config.Database.SSLMode)

	// Also check for DB_PORT as environment variable
	if portStr := os.Getenv("DB_PORT"); portStr != "" {
		var port int
		if _, err := fmt.Sscanf(portStr, "%d", &port); err == nil {
			config.Database.Port = port
		}
	}

	return &config, nil
}

// getEnv returns environment variable value or default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.User,
		c.Database.Password,
		c.Database.SSLMode,
	)
}

// Validate checks if the configuration has all required fields
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port == 0 {
		return fmt.Errorf("database port is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	return nil
}
