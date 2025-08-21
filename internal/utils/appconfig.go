package utils

import (
	"os"
	"path/filepath"

	"github.com/YogeshUpdhyay/ypoker/internal/constants"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Port    string `yaml:"port"`
	Name    string `yaml:"server_name"`
	Version string `yaml:"version"`
}

func DefaultAppConfig() AppConfig {
	return AppConfig{
		Port:    constants.ServerPortDefault,
		Name:    constants.ServerNameDefault,
		Version: constants.ServerVersionDefault,
	}
}

// WriteConfigYML writes the default config to the default path from constants
func WriteAppConfig() error {
	config := DefaultAppConfig()
	dir := constants.ApplicationDataDir
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, constants.ApplicationConfigFileName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	return encoder.Encode(config)
}

// ReadConfigYML reads the config from the default path in constants

func GetAppConfig() (AppConfig, error) {
	var config AppConfig
	path := filepath.Join(constants.ApplicationDataDir, constants.ApplicationConfigFileName)
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultAppConfig(), nil
		}
		return config, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
