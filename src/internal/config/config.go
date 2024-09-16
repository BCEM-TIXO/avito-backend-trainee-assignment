package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
	Storage       StorageConfig
}

type StorageConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Database string `env:"POSTGRES_DATABASE" env-default:"postgres"`
	Username string `env:"POSTGRES_USERNAME" env-default:"postgres"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"postgres"`
	Conn     string `env:"POSTGRES_CONN" env-default:""`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadEnv(instance); err != nil {
			return
		}
	})
	return instance
}
