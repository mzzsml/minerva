package storage

import (
    "context"
    "log"
    "fmt"

    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "github.com/mzzsml/nparse"
)

var ctx context.Context = context.Background()

// Cnnect() handles the authentication and connection to the Postgres db.
// Returns a *pgxpool.Pool object.
func NewConnectionPool(dbfile string) (*sql.DB, error) {
    var dsn string
    dsn = fmt.Sprintf("file:%s?mode=rw", dbfile)
    pool, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return nil, err
    }
    err = pool.Ping()
    if err != nil {
        return nil, err
    }
    return pool, nil
}

// Db holds the connecton pool.
// Using this struct it's possible to query the db without the need of
// reautentication every time.
type DB struct {
    Pool *sql.DB
}

// NewDb creates a new Db instance.
func NewDb(pool *sql.DB) *DB {
    return &DB{pool}
}

// BEFORE INSERTING PORTS, CHECK WETHER THEY ALREADY EXIST
func (d DB) isPortAlreadyExisting() {}

// InsertPorts queries the db and insert the provided ports for the provided
// host.
// It takes the host id and a nparse.Port as input.
func (d DB) InsertPorts(hostId int, p nparse.Port) {
    q := `INSERT INTO
            port (host_id, protocol, port_num, state, reason, service, product, version, extrainfo)
            values ($1, $2, $3, $4, $5, $6, $7, $8, $9);`

    _ = d.Pool.QueryRowContext(
        ctx,
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

func (d DB) IsHostAlreadyExisting(ipv4 string) bool {
    var exists bool

    q := `SELECT EXISTS(SELECT id FROM hosts WHERE ipv4 = $1);`

    err := d.Pool.QueryRowContext(ctx, q, ipv4).Scan(&exists)
    if err == nil && exists {
        return true
    }
    return false
}

func (d DB) InsertHosts(n *nparse.NmapScan) error {
    var id int

    q := `INSERT INTO hosts (ipv4, mac, vendor, hostname) VALUES ($1, $2, $3, $4) RETURNING id;`

    // tmphost temporary holds the information of an host, in simpler form, in
    // ordert to then add it into the database.
    type tmphost struct {
        ipv4     string
        mac      string
        vendor   string
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
            err := d.Pool.QueryRowContext(
                ctx,
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
func (d DB) GetPortsByHostId(hostId int) []nparse.Port {
    var (
        p     nparse.Port
        ports []nparse.Port
    )

    q := `SELECT
            protocol,
            port_num,
            state,
            reason,
            service,
            COALESCE(version, ''), -- if version is null, return an empty string https://sqlite.org/lang_corefunc.html#coalesce
            COALESCE(extrainfo, '')
          FROM
            ports
          WHERE host_id = $1;`

    rows, err := d.Pool.Query(q, hostId)
    if err != nil {
        log.Printf("error in reading ports from DB: %s\n", err)
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(
            &p.Protocol,
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

type host struct {
    Id   int    `json:"id"`
    Ipv4 string `json:"ipv4"`
}

// GetHosts queries the `host` table and returns a list of Host (which is
// different than nparse.Host, and it just holds the hosts' ids and ipv4 addresses).
func (d DB) GetHosts() ([]host, error) {
    var h host
    var hosts []host

    q := `SELECT id, ipv4 FROM hosts;`
    rows, err := d.Pool.Query(q)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(&h.Id, &h.Ipv4)
        if err != nil {
            return nil, err
        }
        hosts = append(hosts, h)
    }
    return hosts, nil
}

func (d DB) GetHostInfoByIPv4(addr string) (nparse.Host, error) {
    var (
        h        nparse.Host
        ports    []nparse.Port
        id       int
        mac      string
        vendor   string
        hostname string
    )
    q := `SELECT id, mac, vendor, hostname FROM hosts WHERE ipv4=$1;`

    err := d.Pool.QueryRowContext(ctx, q, addr).Scan(&id, &mac, &vendor, &hostname)
    if err != nil {
        return h, err
    }

    ports = d.GetPortsByHostId(id)

    addrIPv4 := nparse.Address{
        Addr:     addr,
        AddrType: "ipv4",
    }
    addrMac := nparse.Address{
        Addr:     mac,
        AddrType: "mac",
        Vendor:   vendor,
    }
    addrs := []nparse.Address{addrIPv4, addrMac}

    hostnames := []string{hostname}

    h = nparse.Host{
        Addrs:     addrs,
        Hostnames: hostnames,
        Ports:     ports,
    }
    return h, nil
}
