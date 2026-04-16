package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Mode string

const (
	ModeWarn     Mode = "warn"
	ModeRedirect Mode = "redirect"
)

type Config struct {
	Mode Mode `mapstructure:"mode"`
}

var Default = Config{
	Mode: ModeWarn,
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "pmguard"), nil
}

func Load() (Config, error) {
	dir, err := configPath()
	if err != nil {
		return Default, nil
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(dir)

	viper.SetDefault("mode", string(ModeWarn))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Default, nil // no config yet, use defaults
		}
		return Default, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Default, err
	}
	return cfg, nil
}

func Save(cfg Config) error {
	dir, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	viper.Set("mode", string(cfg.Mode))
	return viper.WriteConfigAs(filepath.Join(dir, "config.yaml"))
}

func (c Config) String() string {
	return fmt.Sprintf("mode: %s", c.Mode)
}
