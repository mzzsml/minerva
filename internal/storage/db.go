package storage

import (
    "log"
    "context"

    "github.com/jackc/pgx/v5/pgxpool"

    "github.com/mzzsml/minerva/internal/parser"
)

func ConnectToDb() *pgxpool.Pool {
    //TODO: ora connection string e' hardcodata.
    // in fututo passarla come file di configurazione (o variabile d'ambiente)
    connStr := "user=samuele password=asd host=192.168.1.48 dbname=minerva sslmode=disable"
    // https://pkg.go.dev/context#Background
    pool, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        log.Fatal(err)
    }
    return pool
}

//func GetHostDetails(id int) {
//    pool := ConnectToDb()
//    defer pool.Close()
//
//    rows, err := pool.Query(context.Background(), "SELECT port_id, state, service FROM port WHERE host_id=1")
//    if err != nil {
//        log.Fatal(err)
//    }
//
//    for rows.Next() {
//        var (
//            port int
//            state string
//            service string
//        )
//        err := rows.Scan(&port, &state, &service)
//        if err != nil {
//            log.Fatal(err)
//        }
//        log.Printf("port %d, state %s, service %s\n", port, state, service)
//    }
//}

func NewHost(h parser.Host) {
    // inserisco host in db.
    // ritorno qualcosa?

    // io devo inserire in due tabelle diverse: host e port
    // prima inserisco in host, ottengo un eventuale id che mi viene ritornato da postgres,
    // uso questo id per inserire le porte?
    pool := ConnectToDb()
    defer pool.Close()

    // prima inseriamo l'ip nella tabella `host`, poi, una volta ottenuto l'id, facciamo l'inserimento delle porte
    var hostId int
    row := pool.QueryRow(context.Background(), "INSERT INTO host (address) VALUES ($1) RETURNING id", h.Addr)
    err := row.Scan(&hostId)
    if err == nil {
        log.Printf("DEBUG: NewHost -- %d", hostId)
    } else {
        log.Fatal(err)
    }

    //inseriamo le porte
    //var p parser.Port
    for _, port := range h.Ports {
        res, err := pool.Exec(context.Background(), "INSERT INTO port (host_id, port_num, state, reason, service, version, extrainfo) VALUES ($1, $2, $3, $4, $5, $6, $7)", hostId, port.PortId, port.State, port.Reason, port.Service, port.Version, port.Extrainfo)
        if err == nil {
            log.Printf("DEBUG: NewHost2 -- %s", res)
        } else {
            log.Fatal(err)
        }
    }
}
