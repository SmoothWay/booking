package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/caarlos0/env"
	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	HTTP  `yaml:"http" json:"http"`
	DB    `yaml:"db" json:"db"`
	Redis `yaml:"redis" json:"redis"`
	Kafka `yaml:"kafka" json:"kafka"`
}

type HTTP struct {
	Port int `env:"HTTP_PORT" yaml:"port" json:"port"`
}
type DB struct {
	Host         string `env:"DB_HOST" yaml:"host" json:"host"`
	Port         int    `env:"DB_PORT" yaml:"port" json:"port"`
	User         string `env:"DB_USER" yaml:"user" json:"user"`
	Password     string `env:"DB_PASSWORD" yaml:"password" json:"password"`
	Name         string `env:"DB_NAME" yaml:"name" json:"name"`
	MigrationDir string `env:"DB_MIGRATION_DIR" yaml:"migration_dir" json:"migration_dir"`
}

type Redis struct {
	Addr     string `env:"REDIS_ADDR" yaml:"addr" json:"addr"`
	Password string `env:"REDIS_PASSWORD" yaml:"password" json:"password"`
	DB       int    `env:"REDIS_DB" yaml:"db" json:"db"`
}

type Kafka struct {
	Brokers []string `env:"KAFKA_BROKERS" yaml:"brokers" json:"brokers"`
	GroupID string   `env:"KAFKA_GROUP_ID" yaml:"group_id" json:"group_id"`
}

func Load(filepath string) (*Config, error) {
	var cfg *Config

	if filepath != "" {
		content, err := os.ReadFile(filepath)
		if err != nil {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}

		switch {
		case hasExt(filepath, ".yaml", ".yml"):
			if err := yaml.Unmarshal(content, &cfg); err != nil {
				return nil, fmt.Errorf("error parsing YAML: %w", err)
			}
		case hasExt(filepath, ".json"):
			if err := json.Unmarshal(content, &cfg); err != nil {
				return nil, fmt.Errorf("error parsing JSON: %w", err)
			}
		default:
		}
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing env vars: %w", err)
	}

	return cfg, nil
}

func hasExt(fileName string, exts ...string) bool {
	for _, ext := range exts {
		if len(fileName) >= len(ext) && fileName[len(fileName)-len(ext):] == ext {
			return true
		}
	}
	return false
}
