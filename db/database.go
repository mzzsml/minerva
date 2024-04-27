package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func CreateNewPool() (*pgxpool.Pool, error) {
	godotenv.Load()
	var dbConnectionString string = os.Getenv("DATABASE_CONNECTION_STRING")
	dbpool, err := pgxpool.New(context.Background(), dbConnectionString)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}
