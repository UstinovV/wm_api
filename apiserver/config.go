package apiserver

import "github.com/UstinovV/wm_api/database"

type Config struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
	Store string `yaml:"database_url"`
	Database *database.Config
}

func NewConfig() *Config {
	return &Config{
		Database:	database.NewConfig(),
	}
}