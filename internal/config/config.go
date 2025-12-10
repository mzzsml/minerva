package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ListenAddress  string `toml:"ListenAddress"`
	DataSourceName string `toml:"DataSourceName"`
}

func ReturnConfigFilePath() (string, error) {
	xdgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(xdgDir, "minerva", "minerva.toml")
	return configFilePath, err
}

// LoadFromFile reads the confiuration parameters from the given
// configuration file.
// It's a method with a pointer receiver, because it needs to modify the state
// of the Config struct.
// It returns an error if the file not exists or if there are problems
// unmarshaling.
func (c *Config) LoadFromFile(filename string) error {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return err
	}

	// NOTE: if file is empty, LoadFromFile as of now doesn't print an error.
	fileContent, err := os.ReadFile(filename)
	err = toml.Unmarshal(fileContent, &c)
	if err != nil {
		return err
	}

	return nil
}

// FindConfigFile simply looks for Minerva configuration files in
// $XDG_CONFIG_HOME and &HOME/.config.
// If none are found it returns a ErrNotExists.
func FindConfigFile() (configFilePath string, err error) {
	configFilePath, err = ReturnConfigFilePath()
	if err != nil {
		return "", err
	}

	_, err = os.Stat(configFilePath)
	// If the file does not exists, return a ErrNotExist error.
	if err != nil && os.IsNotExist(err) {
		return "", err
	}
	// Otherwise, if the file exists, fileInfo will be populated.
	return configFilePath, err
}
