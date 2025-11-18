package cli

import (
    "flag"
    "os"
    "log"

    "github.com/mzzsml/minerva/internal/http"
    "github.com/mzzsml/minerva/internal/storage"
    "github.com/mzzsml/minerva/internal/config"
)

const (
    listenAddrFlagHelp = "listen address."
    dbConnStrFlagHelp  = "databse connection string."
)

func ParseFlags() {
    var (
        listenAddrFlag string
        dbConnStrFlag  string

        conf config.Config
    )

    flag.StringVar(&listenAddrFlag, "listen", "", listenAddrFlagHelp)
    flag.StringVar(&dbConnStrFlag, "db", "", dbConnStrFlagHelp)
    flag.Parse()

    // For config files:
    // If no flags are provided, then parse the config file.
    // Flags have the priority on config files.

    // First, if no flags are passed, use the config file.
    if flag.Parsed() && flag.NFlag() == 0 {
        configFile, err := config.FindConfigFile()
        if err != nil && os.IsNotExist(err) {
            defaultConfig := config.NewDefaultConfig()
            defaultConfig.CreateConfigFile()

            // At this point the file is created and populated.
            // We can assign the default values to the conf struct.
            conf.DataSourceName = defaultConfig.DataSourceName
            conf.ListenAddress = defaultConfig.ListenAddress
        } else if configFile != "" { // Parse the file.
             err := conf.LoadFromFile(configFile)
             if err != nil {
                 log.Fatalf("error: %s\n", err)
             }
         }
    } else if flag.NFlag() > 0 {
        conf.DataSourceName = dbConnStrFlag
        conf.ListenAddress = listenAddrFlag
    }

    var (
        db     *storage.Db
        server http.Server
    )
    pool := storage.NewConnectionPool(conf.DataSourceName)
    defer pool.Close()
    db = storage.NewDb(pool)

    // The server struct is needed to pass the database access to the handlers
    // later.
    server = http.Server{
        Addr:  conf.ListenAddress,
        Store: db,
    }
    startDeamon(server)
}
