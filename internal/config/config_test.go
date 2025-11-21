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

//// Testare creazione directory duplicata.
//func TestMkdirDuplicate(t *testing.T) {
//    testDir := "/tmp/minerva-duplication-test"
//
//    err := Mkdir(testDir)
//    if err != nil {
//        t.Errorf("TestMkdirDuplicate: Mkdir did not create test directory (%s), got error %v", testDir, err)
//    }
//
//    // Now let's try to create the same directory a second time.
//    // If Mkdir() does not throw an error, this test is considered a fail.
//    err = Mkdir(testDir)
//    if err == nil {
//        t.Errorf("TestMkdirDuplicate: Mkdir did not thew error on folder duplication (%s)", testDir)
//    }
//}

func TestWriteToFile(t *testing.T) {
    var c Config

    c = Config{
        ListenAddress: "asd",
        DataSourceName: "qwe",
    }

    err := c.WriteToFile("/tmp/test-minerva/test-file-write")

    if err != nil {
        t.Errorf("TestWriteToFile: could not create te file: %s", err)
    }
}

func TestLoadFromFile (t *testing.T) {
    var c Config
    
    err := c.LoadFromFile("/tmp/test-minerva/test-file-read.toml")

    if err != nil {
        t.Errorf("TestLoadFromFile: error in reading file (%s) (%s)", "/tmp/test-minerva/test-file-read.toml", err)
    }
}

func TestLoadFromFileNonExisting (t *testing.T) {
    var c Config
    
    err := c.LoadFromFile("/tmp/test-minerva/test-file-read-non-existing.toml")

    if err == nil {
        t.Errorf("TestLoadFromFile: error is nil, and it shouldn't (%s)", err)
    }
}

func TestLoadFromFileEmpty (t *testing.T) {
    var c Config
    
    err := c.LoadFromFile("/tmp/test-minerva/test-file-read-empty.toml")

    if err == nil {
        t.Errorf("error is nil, and it shouldn't (%s)", err)
    }
}

func TestFindConfigFile(t *testing.T) {}

func TestCreateConfigFile(t *testing.T) {}
