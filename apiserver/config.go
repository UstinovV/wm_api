package apiserver

type Config struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
	DatabaseUrl string `yaml:"database_ulr"`
}

func NewConfig() *Config {
	return &Config{}
}