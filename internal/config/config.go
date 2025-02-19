package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env          string       `yaml:"env" env-default:"local"`
	JWTSecret   string        `yaml:"jwt_secret" env-default:"your-secret-key"`
	HTTPServer   HTTPServer   `yaml:"http_server"`
	RedisServer  RedisServer  `yaml:"redis_server"`
	PostgresData PostgresData `yaml:"postgres_data"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type RedisServer struct {
	Address  string `yaml:"address" env-default:":8080"`
	Password string `yaml:"password" env-default:""`
	DB       int    `yaml:"db" env-default:"0"`
}

type PostgresData struct {
	Address string `yaml:"address" env-default:"localhost:5432"`
	Name    string `yaml:"name" env-default:"postgres"`
}

func MustLoadConfig() *Config {
	os.Setenv("CONFIG_PATH", "config/local.yaml")
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return &cfg
}
