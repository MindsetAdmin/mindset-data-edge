// internal/config/config.go
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Site      SiteConfig      `yaml:"site"`
	OpcUA     OpcUAConfig     `yaml:"opcua"`
	Discovery DiscoveryConfig `yaml:"discovery"`
	Cloud     CloudConfig     `yaml:"cloud"`
	Cost      CostConfig      `yaml:"cost"`
}

type SiteConfig struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

type OpcUAConfig struct {
	Endpoint       string `yaml:"endpoint"`
	SecurityMode   string `yaml:"security_mode"`
	SecurityPolicy string `yaml:"security_policy"`
}

type DiscoveryConfig struct {
	ScanTimeoutSeconds int `yaml:"scan_timeout_seconds"`
	NoiseFilterHz      int `yaml:"noise_filter_hz"`
	ObservationSeconds int `yaml:"observation_seconds"`
}

type CloudConfig struct {
	Endpoint            string `yaml:"endpoint"`
	APIKey              string `yaml:"api_key"`
	PushIntervalSeconds int    `yaml:"push_interval_seconds"`
}
type CostConfig struct {
	HourlyCost      float64 `yaml:"hourly_cost"`
	TheoreticalRate float64 `yaml:"theoretical_rate"`
	ProductMargin   float64 `yaml:"product_margin"`
	EnergyPrice     float64 `yaml:"energy_price"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
