package config

import (
    "testing"
    "os"
)

func TestReturnConfigFilePath(t *testing.T) {
    want := "/home/samuele/.config/minerva/minerva.toml"

    fp, err := ReturnConfigFilePath()

    if want != fp || err != nil {
        t.Errorf("ReturnConfigFilePath() returns %s, %v; want %s, nil", fp, err, want)
    }
}

func TestNewDefaultConfig(t *testing.T) {
    want := Config{
        ListenAddress: ":8080",
        DataSourceName: "user=samuele password=asd host=192.168.1.48 dbname=minerva sslmode=disable",
    }
    df := NewDefaultConfig()

    if df != want {
        t.Errorf("NewDEfaultConfig() returns %v; want %v", df, want)
    }
}

func TestMkdir(t *testing.T) {
    testDir := "/tmp/minerva-test"
    err := Mkdir(testDir)
    dirCreated, _ := os.Stat(testDir)
    if dirCreated == nil || err != nil {
        t.Errorf("Mkdir did not create test directory (%s), got error %v", testDir, err)
    }
}
