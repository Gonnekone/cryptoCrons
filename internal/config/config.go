package config

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string        `yaml:"env" env:"ENV" env-required:"true"`
	Storage    Storage       `yaml:"storage"`
	HTTPServer HTTPServer    `yaml:"http_server"`
	Interval   time.Duration `yaml:"interval" env-required:"true"`
}

type Storage struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Database string `yaml:"database" env-default:"postgres"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"ADDRESS" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"60s"`
}

func MustLoad() *Config {
	timeout := getEnv("TIMEOUT")
	to, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatal(err)
	}

	idleTimeout := getEnv("IDLE_TIMEOUT")
	it, err := time.ParseDuration(idleTimeout)
	if err != nil {
		log.Fatal(err)
	}

	interval := getEnv("INTERVAL")
	iv, err := time.ParseDuration(interval)
	if err != nil {
		log.Fatal(err)
	}

	cfg := Config{
		Env: getEnv("ENV"),
		Storage: Storage{
			Host:     getEnv("POSTGRES_HOST"),
			Port:     getEnv("POSTGRES_PORT"),
			Database: getEnv("POSTGRES_DB"),
			User:     getEnv("POSTGRES_USER"),
			Password: getEnv("POSTGRES_PASSWORD"),
		},
		HTTPServer: HTTPServer{
			Address:     getEnv("ADDRESS"),
			Timeout:     to,
			IdleTimeout: it,
		},
		Interval: iv,
	}

	return &cfg
}

func (s *Storage) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", s.User, s.Password, s.Host, s.Port, s.Database)
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("error: %s variable not found", key)
	}
	return value
}
