package storage

import (
    "testing"

    "github.com/mzzsml/nparse"
)

const (
    testdb = "testing/test.db"
)

// TODO: creare un db ogni volta che si lanciano i test.
func init() {}

func TestNewConnectionPoolFileNotExist(t *testing.T) {
    pool, err := NewConnectionPool("testing/NotExisting.db")
    if pool != nil {
        t.Errorf("error: pool object should be nil, got %v", pool) 
    }
    if err == nil {
        t.Errorf("error: error should not be nil") 
    }
}

func TestNewConnectionPool(t *testing.T) {
    pool, err := NewConnectionPool(testdb)
    if pool == nil {
        t.Errorf("error: pool object should not be nil") 
    }
    if err != nil {
        t.Errorf("error: error should be nil, got: %s", err) 
    }
}

func TestNewDb(t *testing.T) {
    pool, _ := NewConnectionPool(testdb)
    db := NewDb(pool)
    if db == nil {
        t.Errorf("error: storage.DB object has no been created")
    }
}
