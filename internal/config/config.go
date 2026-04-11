// Package config handles loading and resolving connection configuration for eventctl.
//
// Config is read from ~/.eventctl/config.yaml and allows users to define
// named broker connections so they do not need to repeat addresses on every command.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const configDir = ".eventctl"
const configFile = "config.yaml"

type ConnectionConfig struct {
	Brokers []string `yaml:"brokers,omitempty"`
}

type Config struct {
	Connections map[string]ConnectionConfig `yaml:"connections"`
	Default     string                      `yaml:"default"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("config: cannot resolve home directory: %w", err)
	}
	return filepath.Join(home, configDir, configFile), nil
}

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("config: cannot read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: cannot parse %s: %w", path, err)
	}

	return &cfg, nil
}

func ResolveConnection(cfg *Config, name string) (*ConnectionConfig, error) {

	if name == "" {
		if cfg.Default == "" {
			return nil, errors.New(
				"config: no connection specified and no default connection defined; " +
					"use --connection or set 'default' in ~/.eventctl/config.yaml",
			)
		}
		name = cfg.Default
	}

	conn, ok := cfg.Connections[name]
	if !ok {
		return nil, fmt.Errorf(
			"config: connection %q not found; available connections: [%s]",
			name, joinKeys(cfg.Connections),
		)
	}

	return &conn, nil
}

// BuildConnectionConfig resolves the final ConnectionConfig to use for a command,
// applying overrides from CLI flags on top of any config-file values.
//
// Priority (highest to lowest):
//  1. CLI flags (brokerOverride)
//  2. Named connection from config file (connName)
//  3. Default connection from config file
func BuildConnectionConfig(
	cfg *Config,
	connName string,
	brokerOverride string,
) (*ConnectionConfig, error) {
	var base ConnectionConfig

	if len(cfg.Connections) > 0 || cfg.Default != "" {
		resolved, err := ResolveConnection(cfg, connName)
		if err != nil && brokerOverride == "" {
			return nil, err
		}
		if resolved != nil {
			base = *resolved
		}
	}

	if brokerOverride != "" {
		base.Brokers = splitBrokers(brokerOverride)
	}

	if len(base.Brokers) == 0 {
		return nil, errors.New(
			"config: no broker addresses; provide --broker host:port",
		)
	}

	return &base, nil
}

func splitBrokers(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func joinKeys(m map[string]ConnectionConfig) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}
