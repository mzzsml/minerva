package storage

import (
    "context"
    "log"
    "net"

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

// Db holds the connecton pool.
// Using this struct it's possible to query the db without the need of
// reautentication every time.
type Db struct {
    Pool *pgxpool.Pool
}

// NewDb creates a new Db instance.
func NewDb(pool *pgxpool.Pool) *Db {
    return &Db{pool}
}

// BEFORE INSERTING PORTS, CHECK WETHER THEY ALREADY EXIST
func (d Db) isPortAlreadyExisting() {}

// InsertPorts queries the db and insert the provided ports for the provided
// host.
// It takes the host id and a nparse.Port as input.
func (d Db) InsertPorts(hostId int, p nparse.Port) {
    q := `INSERT INTO
            port (host_id, protocol, port_num, state, reason, service, product, version, extrainfo)
            values ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

    _ = d.Pool.QueryRow(
        context.Background(),
        q,
        hostId,
        p.Protocol,
        p.PortId,
        p.State.State,
        p.State.Reason,
        p.Service.Name,
        p.Service.Product,
        p.Service.Version,
        p.Service.Extrainfo,
    )
}

func (d Db) IsHostAlreadyExisting(ipv4 string) bool {
    var id int

    q := `SELECT id FROM host WHERE ipv4 = $1;`

    err := d.Pool.QueryRow(context.Background(), q, ipv4).Scan(&id)
    if err == nil && id != 0 {
        return true
    }
    return false
}

func (d Db) InsertHosts(n *nparse.NmapScan) error {
    var id int

    q := `INSERT INTO host (ipv4, mac, vendor, hostname) VALUES ($1, $2, $3, $4) RETURNING id;`

    // tmphost temporary holds the information of an host, in simpler form, in
    // ordert to then add it into the database.
    type tmphost struct {
        ipv4 string
        mac string
        vendor string
        hostname string
    }

    for _, h := range n.Hosts {
        var th tmphost
        for _, a := range h.Addrs {
            if a.AddrType == "ipv4" {
                th.ipv4 = a.Addr
                continue
            }
            if a.AddrType == "mac" {
                th.mac = a.Addr
                th.vendor = a.Vendor
                continue
            }
        }
        if !d.IsHostAlreadyExisting(th.ipv4) {
            err := d.Pool.QueryRow(
                context.Background(),
                q,
                th.ipv4,
                th.mac,
                th.vendor,
                th.hostname,
            ).Scan(&id)
            if err != nil {
                return err
            }
            for _, p := range h.Ports {
                d.InsertPorts(id, p)
            }
        } else {
            // TODO: ERROR HOST ALREADY EXISTS.
            return nil
        }
    }
    return nil
}

// GetPorts queries the db and gets every port for a provided host.
// It takes the host ID as input and returns and array of nmap.Port.
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

type Host struct {
    Id int
    Ipv4 net.IP
}
// GetHosts queries the `host` table and returns a list of Host (which is
// different than nparse.Host, and it just holds the hosts' ids and ipv4 addresses).
func (d Db) GetHosts() []Host {
    var host Host
    var hosts []Host

    q := `SELECT id, ipv4 FROM host;`
    rows, err := d.Pool.Query(context.Background(), q)
    if err != nil {
        log.Printf("error in retrieving hosts: %s\n", err)
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&host.Id, &host.Ipv4)
        if err != nil {
            log.Printf("error: %s\n", err)
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
