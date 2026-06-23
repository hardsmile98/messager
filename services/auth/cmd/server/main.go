package main

import (
	"auth/internal/config"
	"fmt"
	"log"
)

func main() {
	config, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println(config)
}
