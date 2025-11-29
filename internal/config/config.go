package config

import (
    "os"

    "github.com/BurntSushi/toml"
)

type Config struct {
    ListenAddress  string `toml:"ListenAddress"`
    DataSourceName string `toml:"DataSourceName"`
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
