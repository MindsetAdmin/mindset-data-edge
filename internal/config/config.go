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
	Mqtt      MqttConfig      `yaml:"mqtt"`
	Cloud     CloudConfig     `yaml:"cloud"`
	Cost      CostConfig      `yaml:"cost"`
}

type MqttConfig struct {
	Broker            string `yaml:"broker"`
	RawTopicPrefix    string `yaml:"raw_topic_prefix"`
	EventsTopicPrefix string `yaml:"events_topic_prefix"`
}

type SiteConfig struct {
	Name string `yaml:"name"`
	ID   string `yaml:"id"`
}

type OpcUAConfig struct {
	Endpoint       string `yaml:"endpoint"`
	SecurityMode   string `yaml:"security_mode"`
	SecurityPolicy string `yaml:"security_policy"`
	Username       string `yaml:"username"`
	// SessionTimeoutSec defaults to 60 when unset.
	SessionTimeoutSec int `yaml:"session_timeout_seconds"`
	// AutoConnect controls whether the edge agent opens the OPC-UA session on boot.
	// Defaults to false so the connection is driven from the frontend (cmd/server).
	AutoConnect bool `yaml:"auto_connect"`
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
