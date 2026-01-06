package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	OutletID          string `json:"outletId"`
	GatewayURL        string `json:"gatewayUrl"`
	AgentVersion      string `json:"agentVersion"`
	ReconnectInterval int    `json:"reconnectInterval"`
	HeartbeatInterval int    `json:"heartbeatInterval"`
	PrinterTimeout    int    `json:"printerTimeout"`
	PrinterPort       int    `json:"printerPort"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	err = json.NewDecoder(file).Decode(cfg)
	return cfg, err
}