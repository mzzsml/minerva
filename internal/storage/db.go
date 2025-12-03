package storage

import (
    "context"
    "log"
    //"net"

    "github.com/jackc/pgx/v5/pgxpool"

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

func (d Db) InsertHost(n *nparse.NmapScan) {
    // Cant use parameters in column header, so we need to do a different
    // insetion per address type (ipv4 and mac).
    qi := `INSERT INTO host (ipv4, vendor) VALUES ($1, $2) RETURNING id;`
    qm := `INSERT INTO host (mac, vendor) VALUES ($1, $2) RETURNING id;`
    var id int
    var ids []int

    for _, h := range n.Hosts {
        for _, a := range h.Addrs {
            if a.AddrType == "ipv4" {
                err := d.Pool.QueryRow(
                    context.Background(),
                    qi,
                    a.Addr,
                    a.Vendor,
                ).Scan(&id)
                if err != nil {
                    log.Printf("error in InsertHost(): %s\n", err)
                }
                ids = append(ids, id)
            }
            if a.AddrType == "mac" {
                err := d.Pool.QueryRow(
                    context.Background(),
                    qm,
                    a.Addr,
                    a.Vendor,
                ).Scan(&id)
                if err != nil {
                    log.Printf("error in InsertHost(): %s\n", err)
                }
                ids = append(ids, id)
            }
        }
    }
}

func (d Db) GetPorts(hostId int) []nparse.Port {
    var (
        p nparse.Port
        ports []nparse.Port
    )

    q := `SELECT
            port_num,
            state,
            reason,
            service,
            COALESCE(version, ''),
            COALESCE(extrainfo, '')
          FROM
            port
          WHERE host_id = $1;`

    rows, err := d.Pool.Query(context.Background(), q, hostId)
    if err != nil {
        log.Printf("error in reading ports from DB: %s\n", err)
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(
            &p.PortId,
            &p.State.State,
            &p.State.Reason,
            &p.Service.Name,
            &p.Service.Version,
            &p.Service.Extrainfo,
        )
        if err != nil {
            log.Printf("errror: %s\n", err)
        }
        ports = append(ports, p)
    }
    return ports
}

// GetHosts queries the `host` table and returns a list of hosts strings containing
// all the IPv4 addresses in it.
func (d Db) GetHosts() []nparse.Host {
    var (
        hostId int
        host nparse.Host
        hosts []nparse.Host
    )

    q := `SELECT id FROM host;`
    rows, err := d.Pool.Query(context.Background(), q)
    if err != nil {
        log.Printf("error in retrieving hosts: %s\n", err)
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&hostId)
        if err != nil {
            log.Printf("error: %s\n", err)
        }
        host = nparse.Host{
            Ports: d.GetPorts(hostId),
        }
        hosts = append(hosts, host)
    }
    return hosts
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
