package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	DatabaseURL string
	Port        int
}

var (
	Config Configuration
)

func init() {
	file, err := os.Open("static/config.json")
	if err != nil {
		log.Println("Couldn't find static/config.json")
		os.Exit(1)
	}

	defer file.Close()

	parser := json.NewDecoder(file)
	err = parser.Decode(&Config)
	if err != nil {
		log.Println("Malformed JSON in configuration file")
		os.Exit(1)
	}
}
