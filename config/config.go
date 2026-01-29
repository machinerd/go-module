package config

import (
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	k *koanf.Koanf
}

func New(configPath string) (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, err
	}

	if err := k.Load(env.Provider("ENV_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "ENV_"))
	}), nil); err != nil {
		return nil, err
	}

	return &Config{k: k}, nil
}

func (c *Config) GetConfig() *koanf.Koanf {
	return c.k
}

func (c *Config) Get(key string) interface{} {
	return c.k.Get(key)
}
func (c *Config) String(key string) string {
	return c.k.String(key)
}
func (c *Config) Strings(key string) []string {
	return c.k.Strings(key)
}
func (c *Config) Int(key string) int {
	return c.k.Int(key)
}
func (c *Config) Bool(key string) bool {
	return c.k.Bool(key)
}
func (c *Config) Float64(key string) float64 {
	return c.k.Float64(key)
}
