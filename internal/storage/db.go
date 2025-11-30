package storage

import (
    "context"
    "log"
    "net"

    "github.com/jackc/pgx/v5/pgxpool"

    //"github.com/mzzsml/minerva/internal/model"
    "github.com/mzzsml/nparse"
)

// Cnnect() handles the authentication and connection to the Postgres db.
// Returns a *pgxpool.Pool object.
func NewConnectionPool(dataSourceName string) *pgxpool.Pool {
    // https://pkg.go.dev/context#Background
    pool, err := pgxpool.New(context.Background(), dataSourceName)
    if err != nil {
        log.Fatalf("ERROR: ConnectToDb: %s", err)
    }
    return pool
}

type Db struct {
    Pool *pgxpool.Pool // so i can use like storage.db.Query('do stuff')...
}

func NewDb(pool *pgxpool.Pool) *Db {
    return &Db{pool}
}

// altrimenti non compila
func (d Db) InsertHost(n *nparse.NmapScan) {
    q := `INSERT INTO host (address) VALUES ($1) RETURNING id;`
    var id int

    for _, h := range n.Hosts {
        for _, a := range h.Addrs {
            if a.AddrType == "ipv4" {
                err := d.Pool.QueryRow(context.Background(), q, a.Addr).Scan(&id)
                if err != nil {
                    log.Printf("error in Scan(): %s\n", err)
                }
                log.Printf("%d\n", id)
            }
        }
    }
}

//func (d Db) GetHosts() (hosts []nparse.Host) {
func (d Db) GetHosts() {
    q := `SELECT address FROM host;`
    
    //var h nparse.Host
    rows, err := d.Pool.Query(context.Background(), q)
    if err != nil {
        log.Printf("error in retrieving hosts: %s\n", err)
    }
    defer rows.Close()
    var addr net.IP
    for rows.Next() {
        //err = rows.Scan(&h.Addrs[0].Addr)
        err = rows.Scan(&addr)
        if err != nil {
            log.Printf("error: %s\n", err)
        }
        log.Printf("%s\n", addr.String())
        //hosts = append(hosts, h)
    }
    //return
}

// NewHost handles the insertion of a new host and its details to the db.
//func NewHost(h model.Host) {
//    // First connect to the db.
//    pool := ConnectToDb()
//    // After we're done, close the connection.
//    defer pool.Close()
//
//    var hostId int
//    // Insert the new host into the `host` table.
//    // If the operation is successful, Postgres will return the new host's id.
//    // We're gonna use this id to insert the ports to the `port` table, so that
//    // the entries are linked.
//    // TODO: use prepared statements?
//    row := pool.QueryRow(context.Background(), "INSERT INTO host (address) VALUES ($1) RETURNING id", h.Addr)
//    err := row.Scan(&hostId)
//    if err != nil {
//        // TODO: instead of closing the application, we should send something like http error 500,
//        // and keep this process running
//        log.Fatal(err)
//    }
//
//    // Now insert the ports.
//    // Because ports are more than one, we have to cycle through them
//    // TODO: use prepared statements?
//    for _, port := range h.Ports {
//        _, err := pool.Exec(context.Background(), "INSERT INTO port (host_id, port_num, state, reason, service, version, extrainfo) VALUES ($1, $2, $3, $4, $5, $6, $7)", hostId, port.PortId, port.State, port.Reason, port.Service, port.Version, port.Extrainfo)
//        if err != nil {
//            log.Fatal(err)
//        }
//    }
//}
