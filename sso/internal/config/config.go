package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

const configPathEnv = "config_path"

type Config struct {
	Env    string `yaml:"env" env-default:"prod"`
	Server Server `yaml:"server" env-required:"true"`
}

type Server struct {
	Addr        string        `yaml:"addr" env-default:"localhost:7777"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

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
