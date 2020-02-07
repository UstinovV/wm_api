package main

import (
	"github.com/UstinovV/wm_api/apiserver"
	_"database/sql"
	_"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)


func main() {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	config :=  apiserver.NewConfig()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}