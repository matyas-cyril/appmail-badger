package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	DBName   string `json:"db_name"`
	User     string `json:"user"`
	Password string `json:"password"`
	SSLMode  bool   `json:"sslmode"`
	Timeout  uint8  `json:"timeout"`
}

var DefaultConfig = func() Config {

	return Config{
		Host:    "127.0.0.1",
		Port:    5432,
		SSLMode: false,
		Timeout: 5,
	}

}

func Load(jsonFile string) (*Config, error) {

	file, err := os.Open(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load json file : %v", err)
	}
	defer file.Close()

	cfg := DefaultConfig()
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode Json : %v", err)
	}

	return &cfg, nil
}
