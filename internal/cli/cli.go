package cli

import (
    "flag"

    "github.com/mzzsml/minerva/internal/storage"
    "github.com/mzzsml/minerva/internal/http"
)


const (
    listenAddrFlagHelp = "listen address."
    dbConnStrFlagHelp = "databse connection string."
)

func ParseFlags () {
    var (
        listenAddrFlag string
        dbConnStrFlag string
    )

    flag.StringVar(&listenAddrFlag, "listen", ":8080", listenAddrFlagHelp)
    flag.StringVar(&dbConnStrFlag, "db", "user=samuele password=asd host=192.168.1.48 dbname=minerva sslmode=disable", dbConnStrFlagHelp)

    flag.Parse()

    var (
        db *storage.Db

        server http.Server
    )
    if dbConnStrFlag != "" {
        // creo un pool
        pool := storage.NewConnectionPool(dbConnStrFlag)
        defer pool.Close()

        // creo un oggetto storage.Db, con pool come parametro
        db = storage.NewDb(pool)

        // creo questo struct cosi' posso mettere il db a disposizione degli
        // handler.
        server = http.Server {
            Addr: listenAddrFlag,
            Store: db,
        }
    }
    startDeamon(server)
}
