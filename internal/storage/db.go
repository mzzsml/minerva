package storage

import (
    "context"
    "fmt"

    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "github.com/mzzsml/nparse"
)

var ctx context.Context = context.Background()

// NewConnectionPool opens a connection to the SQLite database, and returns a
// connection pool or an error;
func NewConnectionPool(dbfile string) (*sql.DB, error) {
    // TODO: io qui passo semplicemente il file, senza prima verificare che esso esista.
    pool, err := sql.Open("sqlite3", "file:"+dbfile+"?mode=rw")
    if err != nil {
        return nil, fmt.Errorf("NewConnectionPool: error in opening database file: %s", err)
    }
    if err = pool.Ping(); err != nil {
        return nil, fmt.Errorf("NewConnectionPool: error in connecting to database: %s", err)
    }
    return pool, nil
}

// DB holds the connecton pool.
// Using this struct it's possible to query the db without the need of
// reautentication every time.
type DB struct {
    Pool *sql.DB
}

// NewDb creates a new Db instance.
func NewDb(pool *sql.DB) *DB {
    return &DB{pool}
}

// insertPorts queries the db and insert the provided ports for the provided
// host.
func (d DB) insertPorts(hostId int, ports []nparse.Port) error {
    q := `INSERT INTO ports values ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING rowid;`
    var row int

    for _, port := range ports {
        err := d.Pool.QueryRowContext(
            ctx,
            q,
            hostId,
            port.Protocol,
            port.PortId,
            port.State.State,
            port.State.Reason,
            port.Service.Name,
            port.Service.Product,
            port.Service.Version,
            port.Service.Extrainfo,
        ).Scan(&row)
        if err != nil {
            return fmt.Errorf("insertPorts: error in inserting ports to database: %s", err)
        }
    }
    return nil
}

func (d DB) hostExists(ipv4 string) bool {
    var exists bool
    q := `SELECT EXISTS(SELECT id FROM hosts WHERE ipv4 = $1);`
    err := d.Pool.QueryRowContext(ctx, q, ipv4).Scan(&exists)
    if err == nil && exists {
        return true
    }
    return false
}

func (d DB) insertHostnames(hostId int, hostnames []nparse.Hostname) error {
    var row int
    q := `INSERT INTO
            hostnames (host_id, hostname, type)
            values ($1, $2, $3)
          RETURNING rowid;`
    for _, h := range hostnames {
        err := d.Pool.QueryRowContext(ctx, q, hostId, h.Hostname, h.Type).Scan(&row)
        if err != nil {
            return fmt.Errorf("insertHostnames: error in inserting hostnames to database: %s", err)
        }
    }
    return nil
}

// InsertHosts inserts a new host in the database, if the host does not already
// exists.
func (d DB) InsertHosts(n *nparse.NmapScan) error {
    var id int
    q := `INSERT INTO hosts (ipv4, mac, vendor, status, reason, scanned_at)
            VALUES ($1, $2, $3, $4, $5, $6)
            RETURNING id;`
    for _, h := range n.Hosts {
        ipv4 := h.AddrInfo("ipv4")
        mac := h.AddrInfo("mac")
        if d.hostExists(ipv4.Addr) {
            // TODO: errore, host gia' esiste. Ma se abbiamo tanti host, e solo uno esiste?
            // Non possiamo inviare NotModified e fermare tutto. Lasciamo continue e basta?
            continue
        }
        err := d.Pool.QueryRowContext(
            ctx,
            q,
            ipv4.Addr,
            mac.Addr,
            mac.Vendor,
            h.Status.State,
            h.Status.Reason,
            h.ScannedAt,
        ).Scan(&id)
        if err != nil {
            return fmt.Errorf("InserHosts: error in inserting host(s) to database: %s", err)
        }
        if err := d.insertPorts(id, h.Ports); err != nil {
            return fmt.Errorf("InsertHosts: %s", err)
        }
        if err := d.insertHostnames(id, h.Hostnames); err != nil {
            return fmt.Errorf("InsertHosts: %s", err)
        }
    }
    return nil
}

// GetPorts queries the db and gets every port for a provided host.
// It takes the host ID as input and returns and array of nmap.Port.
func (d DB) GetPortsByHostId(hostId int) ([]nparse.Port, error) {
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
            product,
            COALESCE(version, ''), -- if version is null, return an empty string https://sqlite.org/lang_corefunc.html#coalesce
            COALESCE(extrainfo, '')
          FROM
            ports
          WHERE host_id = $1;`

    rows, err := d.Pool.Query(q, hostId)
    if err != nil {
        return nil, fmt.Errorf("GetPortsByHostId: error in reading ports from database: %s", err)
    }
    defer rows.Close()

    for rows.Next() {
        err := rows.Scan(
            &p.Protocol,
            &p.PortId,
            &p.State.State,
            &p.State.Reason,
            &p.Service.Name,
            &p.Service.Product,
            &p.Service.Version,
            &p.Service.Extrainfo,
        )
        if err != nil {
            return nil, fmt.Errorf("GetPortsByHostID: %s", err)
        }
        ports = append(ports, p)
    }
    return ports, nil
}

func (d DB) getHostnames(hostId int) ([]nparse.Hostname, error) {
    var (
        q = `SELECT
                COALESCE(hostname, ''),
                COALESCE(type, '')
            FROM hostnames 
            WHERE host_id = $1;`
        hostname  nparse.Hostname
        hostnames []nparse.Hostname
    )

    rows, err := d.Pool.Query(q, hostId)
    if err != nil {
        return nil, fmt.Errorf("getHostnames: error in getting hostnames from database: %s", err)
    }
    defer rows.Close()
    for rows.Next() {
        if err := rows.Scan(&hostname.Hostname, &hostname.Type); err != nil {
            return nil, fmt.Errorf("getHostnames: %s", err)
        }
        hostnames = append(hostnames, hostname)
    }
    return hostnames, nil
}

type host struct {
    Id   int    `json:"id"`
    Ipv4 string `json:"ipv4"`
}

// GetHosts queries the `host` table and returns a list of Host (which is
// different than nparse.Host, and it just holds the hosts' ids and ipv4 addresses).
func (d DB) GetHosts() ([]host, error) {
    var (
        h     host
        hosts []host
        q     = `SELECT id, ipv4 FROM hosts;`
    )
    rows, err := d.Pool.Query(q)
    if err != nil {
        return nil, fmt.Errorf("GetHosts: error in getting hosts from database: %s", err)
    }
    defer rows.Close()
    for rows.Next() {
        err = rows.Scan(&h.Id, &h.Ipv4)
        if err != nil {
            return nil, fmt.Errorf("GetHosts: %s", err)
        }
        hosts = append(hosts, h)
    }
    return hosts, nil
}

func (d DB) GetHostInfoByIPv4(addr string) (nparse.Host, error) {
    var (
        id                        int
        mac                       string
        vendor                    string
        status, reason, scannedat string
        q                         = `SELECT
                id,
                COALESCE(mac, ''),
                COALESCE(vendor, ''),
                status,
                reason,
                scanned_at
            FROM hosts
            WHERE ipv4=$1;`
    )
    err := d.Pool.QueryRowContext(ctx, q, addr).Scan(&id, &mac, &vendor, &status, &reason, &scannedat)
    if err != nil {
        return nparse.Host{}, fmt.Errorf("GetHostInfoByIPv4: error in getting host info: %s", err)
    }
    addrIPv4 := nparse.Address{Addr: addr, AddrType: "ipv4"}
    addrMac := nparse.Address{mac, "mac", vendor}
    hostnames, err := d.getHostnames(id)
    if err != nil {
        return nparse.Host{}, fmt.Errorf("GetHostInfoByIPv4: %s", err)
    }
    ports, err := d.GetPortsByHostId(id)
    if err != nil {
        return nparse.Host{}, fmt.Errorf("GetHostInfoByIPv4: %s", err)
    }
    h := nparse.Host{
        Status:    nparse.State{status, reason},
        Addrs:     []nparse.Address{addrIPv4, addrMac},
        Hostnames: hostnames,
        Ports:     ports,
        //Extraports
        ScannedAt: scannedat,
    }
    return h, nil
}

func (d DB) DeleteHost(addr string) bool {
    q := `DELETE FROM hosts WHERE ipv4 = $1;`
    if !d.hostExists(addr) {
        return false
    }
    if _, err := d.Pool.ExecContext(ctx, q, addr); err != nil {
        return false
    }
    return true
}
