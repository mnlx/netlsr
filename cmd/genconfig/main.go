package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mnlx/netlsr/internal/config"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run cmd/genconfig/main.go <output_dir>")
		fmt.Println("Example: go run cmd/genconfig/main.go ./configs")
		os.Exit(1)
	}

	outputDir := os.Args[1]

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Generate default config
	defaultConfig := config.DefaultConfig()

	// Save as config.yaml
	configPath := filepath.Join(outputDir, "config.yaml")
	if err := config.SaveConfig(defaultConfig, configPath); err != nil {
		log.Fatalf("Failed to save config: %v", err)
	}

	fmt.Printf("Default configuration saved to: %s\n", configPath)
	fmt.Println("You can now customize the values in this file.")
}
