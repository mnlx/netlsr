package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application
type Config struct {
	Client ClientConfig `yaml:"client"`
	Server ServerConfig `yaml:"server"`
	Debug  bool         `yaml:"debug" env:"NETLSR_DEBUG"`
}

// ClientConfig holds client-specific configuration
type ClientConfig struct {
	TunName    string `yaml:"tun_name" env:"NETLSR_CLIENT_TUN_NAME"`
	TunCIDR    string `yaml:"tun_cidr" env:"NETLSR_CLIENT_TUN_CIDR"`
	LocalIP    string `yaml:"local_ip" env:"NETLSR_CLIENT_LOCAL_IP"`
	PeerIP     string `yaml:"peer_ip" env:"NETLSR_CLIENT_PEER_IP"`
	ServerAddr string `yaml:"server_addr" env:"NETLSR_CLIENT_SERVER_ADDR"`
	Port       int    `yaml:"port" env:"NETLSR_CLIENT_PORT"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	TunName  string `yaml:"tun_name" env:"NETLSR_SERVER_TUN_NAME"`
	TunCIDR  string `yaml:"tun_cidr" env:"NETLSR_SERVER_TUN_CIDR"`
	LocalIP  string `yaml:"local_ip" env:"NETLSR_SERVER_LOCAL_IP"`
	PeerIP   string `yaml:"peer_ip" env:"NETLSR_SERVER_PEER_IP"`
	Port     int    `yaml:"port" env:"NETLSR_SERVER_PORT"`
	ExtIface string `yaml:"ext_iface" env:"NETLSR_SERVER_EXT_IFACE"`
	Debug    bool   `yaml:"debug" env:"NETLSR_SERVER_DEBUG"`
}

// DefaultConfig returns a config with default values
func DefaultConfig() *Config {
	return &Config{
		Client: ClientConfig{
			TunName:    "tun0",
			LocalIP:    "10.0.0.1",
			PeerIP:     "10.0.0.2",
			ServerAddr: "",
			Port:       8080,
		},
		Server: ServerConfig{
			TunName:  "tun0",
			LocalIP:  "10.0.0.1",
			PeerIP:   "10.0.0.2",
			Port:     8080,
			ExtIface: "",
			Debug:    false,
		},
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()

	// Load from file if provided
	if configPath != "" {
		if err := loadFromFile(configPath, config); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
		if config.Debug {
			log.Println(config)
		}
	}

	// Override with environment variables
	loadFromEnv(config)

	return config, nil
}

// loadFromFile loads configuration from a YAML file
func loadFromFile(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, config)
}

// loadFromEnv overrides config values with environment variables
func loadFromEnv(config *Config) {
	// Client config
	if val := os.Getenv("NETLSR_CLIENT_TUN_NAME"); val != "" {
		config.Client.TunName = val
	}
	if val := os.Getenv("NETLSR_CLIENT_LOCAL_IP"); val != "" {
		config.Client.LocalIP = val
	}
	if val := os.Getenv("NETLSR_CLIENT_PEER_IP"); val != "" {
		config.Client.PeerIP = val
	}
	if val := os.Getenv("NETLSR_CLIENT_SERVER_ADDR"); val != "" {
		config.Client.ServerAddr = val
	}
	if val := os.Getenv("NETLSR_CLIENT_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.Client.Port = port
		}
	}

	// Server config
	if val := os.Getenv("NETLSR_SERVER_TUN_NAME"); val != "" {
		config.Server.TunName = val
	}
	if val := os.Getenv("NETLSR_SERVER_LOCAL_IP"); val != "" {
		config.Server.LocalIP = val
	}
	if val := os.Getenv("NETLSR_SERVER_PEER_IP"); val != "" {
		config.Server.PeerIP = val
	}
	if val := os.Getenv("NETLSR_SERVER_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			config.Server.Port = port
		}
	}
	if val := os.Getenv("NETLSR_SERVER_EXT_IFACE"); val != "" {
		config.Server.ExtIface = val
	}
	if val := os.Getenv("NETLSR_SERVER_DEBUG"); val != "" {
		config.Server.Debug = strings.ToLower(val) == "true"
	}
}

// SaveConfig saves the current configuration to a YAML file
func SaveConfig(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
