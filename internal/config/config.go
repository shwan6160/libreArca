package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	WikiName string `yaml:"wiki_name"`
	BbsName  string `yaml:"bbs_name"`
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	type rawConfig Config
	var raw rawConfig
	if err := value.Decode(&raw); err != nil {
		return err
	}
	*c = Config(raw)
	if c.BbsName == "" {
		c.BbsName = c.WikiName
	}
	return nil
}

var AppConfig Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &AppConfig)
}
