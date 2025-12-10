package cli

import (
	"flag"
	"os"

	"github.com/mzzsml/minerva/internal/config"
	"github.com/mzzsml/minerva/internal/http"
	"github.com/mzzsml/minerva/internal/storage"
)

const (
	listenAddrFlagHelp = "listen address."
	dbConnStrFlagHelp  = "databse connection string."
)

// ParseFlags parses the command line flags provided.
// If no flags are provided, then it uses the configuration files, located either in
// $XDG_CONFIG_HOME or $HOME/.config/minerva.
//
// Note that flags have the priority on files, so, if they are provided and the files
// already exist, ParseFlags ignores the files and uses the given command line flags.
func ParseFlags() error {
	var (
		listenAddrFlag string
		dbConnStrFlag  string

		conf config.Config

		db     *storage.Db
		server http.Server
	)

	flag.StringVar(&listenAddrFlag, "listen", "", listenAddrFlagHelp)
	flag.StringVar(&dbConnStrFlag, "db", "", dbConnStrFlagHelp)
	flag.Parse()

	// If no flags are given, use the config file.
	if flag.Parsed() && flag.NFlag() == 0 {
		configFile, err := config.FindConfigFile()
		if err != nil && os.IsNotExist(err) {
			// NOTE: for now just return an error.
			return err
		} else if configFile != "" {
			// If the file exists, then parse it.
			err := conf.LoadFromFile(configFile)
			if err != nil {
				return err
			}
		}
	} else if flag.NFlag() > 0 {
		// If flags are given, read their value and put it in conf.
		conf.DataSourceName = dbConnStrFlag
		conf.ListenAddress = listenAddrFlag
	}

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

	return nil
}
