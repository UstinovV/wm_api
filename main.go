package main

import (
	_ "database/sql"
	_ "fmt"
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

	config := NewConfig()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	serv := NewServer(config)

	err = serv.Start()
	if err != nil {
		log.Fatal("Server error", err)
	}

}
