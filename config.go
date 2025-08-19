package main

import (
	"encoding/json"
	// "fmt"
	"io"
	"log"
	"os"
)

// var Name = "yadhu"

func ConfigLoader() *Config{
	file, err := os.Open("config.json")
	if err != nil {
		log.Println("❌Error in file opening", err)
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error in file reading", err)
	}
	var cfg Config
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return &cfg

}
