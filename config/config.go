package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type RuleConfig struct {
	Enabled bool                   `yaml:"enabled"`
	Options map[string]interface{} `yaml:"options"`
}

type Config struct {
	Rules map[string]RuleConfig `yaml:"rules"`
}

func Load(configPath string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Printf("cannot read config %q %v using default rule", configPath, err)
		return DefaultConfig()
	}
	return &cfg
}

func DefaultConfig() *Config {
	return &Config{
		Rules: map[string]RuleConfig{
			"lowercase": {
				Enabled: true,
				Options: map[string]interface{}{},
			},
			"english": {
				Enabled: true,
				Options: map[string]interface{}{},
			},
			"symbols": {
				Enabled: true,
				Options: map[string]interface{}{
					"allowed": []string{".", "-", "_", "/", ":", "=", ",", "%"},
				},
			},
			"sensitive": {
				Enabled: true,
				Options: map[string]interface{}{
					"patterns": []string{
						"(?i)\\b(password|passwd|pwd)\\b",
						"(?i)\\b(api[_-]?key|apikey)\\b",
						"(?i)\\btoken\\b",
					},
				},
			},
		},
	}
}
