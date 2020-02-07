package apiserver

import "github.com/UstinovV/wm_api/database"

type Config struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
	Database *database.Config
}

func NewConfig() *Config {
	return &Config{
		Database:	database.NewConfig(),
	}
}