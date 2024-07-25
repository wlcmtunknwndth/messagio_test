package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const configPathEnv = "config_path"

type Config struct {
	Env      string        `yaml:"env" env-default:"prod"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"5m"`
	Server   Server        `yaml:"server" env-required:"true"`
	DB       DataBase      `yaml:"db"`
}

type Server struct {
	Addr        string        `yaml:"addr" env-default:"localhost:7777"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type DataBase struct {
	DbUser  string `yaml:"user" env-default:"postgres"`
	DbPass  string `yaml:"pass" env-default:"postgres"`
	DbName  string `yaml:"name" env-default:"postgres"`
	SslMode string `yaml:"ssl_mode" env-default:"disable"`
	Port    string `yaml:"port" env-default:"5432"`
}

func MustLoad() *Config {
	const op = "inner.config.MustLoad"

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
