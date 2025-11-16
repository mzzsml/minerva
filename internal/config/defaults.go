package config

import (
    "os"
    "log"

    "github.com/BurntSushi/toml"
)

const (
    defaultListenAddress  = ":8080"
    defaultDataSourceName = "user=samuele password=asd host=192.168.1.48 dbname=minerva sslmode=disable"

    confdir = "./etc/minerva/"
)

type Config struct {
    ListenAddress  string `toml:"ListenAddress"`
    DataSourceName string `toml:"DataSourceName"`
}

func NewDefaultConfig() Config {
    return Config{
        ListenAddress:  defaultListenAddress,
        DataSourceName: defaultDataSourceName,
    }
}

func Mkdir(d string) error {
    _, err := os.Stat(d)
    // la cartella non esiste, ne creo una nuova
    if err != nil && os.IsNotExist(err) {
        log.Printf("INFO -- /etc/minerva not found, creating it")
        err = os.Mkdir(d, 0755)
        if err != nil {
            log.Fatalf("ERROR -- could not create /etc/minerva: %s", err)
            return err
        }
    }
    return err
}

func (c Config) WriteToFile() error {
    // prima creo cartella in /etc/minerva
    // poi metto minerva.toml
    err := Mkdir(confdir)
    if err != nil {
        return err
    }
    // se file minerva.toml non esiste, lo creiamo
    // ora lo creo ogni volta e lo cancello a mano
    filename := confdir + "minerva.toml"
    marshaled, err := toml.Marshal(c)
    err = os.WriteFile(filename, marshaled, 0644)
    if err != nil {
        log.Fatalf("error in writing toml file: %s", err)
    }
    return err
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
