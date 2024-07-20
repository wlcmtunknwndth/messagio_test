package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env    string   `yaml:"env" env-default:"prod"`
	Server Server   `yaml:"server" env-required:"true"`
	DB     DataBase `yaml:"db" env-required:"true"`
	Broker Broker   `yaml:"broker" env-required:"true"`
}

type Server struct {
	Address     string        `yaml:"address" env-default:"localhost:8888"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"10s"`
}

type DataBase struct {
	DbUser  string `yaml:"user" env-default:"postgres"`
	DbPass  string `yaml:"pass" env-default:"postgres"`
	DbName  string `yaml:"name" env-default:"postgres"`
	SslMode string `yaml:"ssl_mode" env-default:"disable"`
	Port    string `yaml:"port" env-default:"5432"`
}

type Broker struct {
}

const configPathEnv = "config_path"

func MustLoad() *Config {
	const op = "internal.config.MustLoad"

	path, ok := os.LookupEnv(configPathEnv)
	if !ok || path == "" {
		panic(op + ": config path is empty")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(op + ": failed to read config:" + err.Error())
	}

	return &cfg
}
