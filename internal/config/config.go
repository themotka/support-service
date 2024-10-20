package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Environment string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage" env-required:"true"`
}

func MustLoad() *Config {
	path := getConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file not found" + path)
	}
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read cfg" + err.Error())
	}
	return &cfg
}

// main.go --config=./path...
// CONFIG_PATH=./path/to/config/file.yaml main.go
func getConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "config path")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}
