package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/minerva/db"
)

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetProjects() string {
	conn, err := db.CreateNewPool()
	if err != nil {
		log.Fatal(err)
	}

	rows, _ := conn.Query(context.Background(), "SELECT id, name FROM project")

	projects, _ := pgx.CollectRows(rows, pgx.RowToStructByName[Project])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		os.Exit(1)
	}

	fmt.Println(projects)

	project, _ := json.Marshal(projects)
	return string(project)
}
