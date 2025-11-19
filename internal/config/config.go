package config

import (
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	defaultListenAddress  = ":8080"
	defaultDataSourceName = "user=samuele password=asd host=192.168.1.48 dbname=minerva sslmode=disable"
)

type Config struct {
	ListenAddress  string `toml:"ListenAddress"`
	DataSourceName string `toml:"DataSourceName"`
}

func returnConfigFilePath() (string, error) {
	xdgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configFilePath := filepath.Join(xdgDir, "minerva", "minerva.toml")
	return configFilePath, err
}

func NewDefaultConfig() Config {
	return Config{
		ListenAddress:  defaultListenAddress,
		DataSourceName: defaultDataSourceName,
	}
}

// Mkdir simply creates a directory, if it does not already exists.
func Mkdir(d string) error {
	_, err := os.Stat(d)
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(d, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteToFile writes the parsed Toml configuration to the configuration file in
// $XDG_CONFIG_HOME.
// It only returns an error.
func (c Config) WriteToFile(configFilePath string) error {
	marshaled, err := toml.Marshal(c)

	// os.Create can't create a file is its parent directory also doesn't exist.
	err = Mkdir(path.Dir(configFilePath))
	if err != nil {
		return err
	}

	// os.WriteFile can't create a file.
	_, err = os.Create(configFilePath)
	if err != nil {
		return err
	}

	// os.Create creates a file with a mask of 0644.
	// For security reason a 0600 is better, since this configuration file is
	// going to contain credentials.
	err = os.Chmod(configFilePath, 0600)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, marshaled, 0600)
	if err != nil {
		return err
	}
	return nil
}

// LoadFromFile leads the confiuration parameters from the provided
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

	fileContent, err := os.ReadFile(filename)
	err = toml.Unmarshal(fileContent, &c)
	if err != nil {
		return err
	}

	return nil
}

// FindConfigFile simply looks for Minerva configuration files in
// $XDG_CONFIG_HOME and &HOME/.config.
// If none are found it return ErrNotExists.
func FindConfigFile() (configFilePath string, err error) {
	configFilePath, err = returnConfigFilePath()
	if err != nil {
		return "", err
	}

	_, err = os.Stat(configFilePath)
	// If the file does not exists, return a ErrNotExist error.
	if err != nil && os.IsNotExist(err) {
		return "", err
	}
	// If file exists, fileInfo will be populated.
	// We can return configFilePath (this was for an "if" check i've deleted)
	return configFilePath, err
}

// CreateConfigFile creates a configuration file in $XDG_CONFIG_HOME, if it's
// not empty. Otherwise it creates it in $HOME/.config/.
// See https://pkg.go.dev/os#UserConfigDir.
func (c Config) CreateConfigFile() error {
	configFilePath, err := returnConfigFilePath()

	_, err = FindConfigFile()
	if err != nil && os.IsNotExist(err) {
		err = c.WriteToFile(configFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}
