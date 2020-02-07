package database

type Config struct {
	DatabaseUrl string `yaml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}