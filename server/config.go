package server

type Config struct {
	Port         string `yaml:"port"`
	Host         string `yaml:"host"`
	DBConnection string `yaml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
