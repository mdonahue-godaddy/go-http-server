package config

import (
	"encoding/json"
	"io/ioutil"
)

// Settings contains values loaded from a config file
type Settings struct {
	Service struct {
		GracefulShutdownDelaySeconds string `json:"gracefulShutdownDelaySeconds" yaml:"gracefulShutdownDelaySeconds" mapstructure:"gracefulShutdownDelaySeconds"`
		HTTP                         struct {
			Server struct {
				IPv4Address string `json:"ipv4address" yaml:"ipv4address" mapstructure:"ipv4address"`
				Port        uint16 `json:"port" yaml:"port" mapstructure:"port"`
			} `json:"server" yaml:"server" mapstructure:"server"`
		} `json:"http" yaml:"http" mapstructure:"http"`
	} `json:"service" yaml:"service" mapstructure:"service"`
	Metrics struct {
		HTTP struct {
			Server struct {
				IPv4Address string `json:"ipv4address" yaml:"ipv4address" mapstructure:"ipv4address"`
				Port        uint16 `json:"port" yaml:"port" mapstructure:"port"`
			} `json:"server" yaml:"server" mapstructure:"server"`
		} `json:"http" yaml:"http" mapstructure:"http"`
		HealthcheckEnabled bool `json:"healthcheckEnabled" yaml:"healthcheckEnabled" mapstructure:"healthcheckEnabled"`
		PPRofEnabled       bool `json:"pprofEnabled" yaml:"pprofEnabled" mapstructure:"pprofEnabled"`
	} `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
	Logging struct {
		Level string `json:"level" yaml:"level" mapstructure:"level"`
	} `json:"logging" yaml:"logging" mapstructure:"logging"`
}

// LoadSettings loads the Settings from JSON file.
func LoadSettings(jsonFile string) (*Settings, error) {
	var settings Settings

	// read file
	appConfigString, rerr := ioutil.ReadFile(jsonFile)
	if rerr != nil {
		return nil, rerr
	}

	jerr := json.Unmarshal([]byte(appConfigString), &settings)
	if jerr != nil {
		return nil, jerr
	}

	return &settings, nil
}
