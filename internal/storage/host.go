package storage

import (
    "context"
    "log"

    //"github.com/jackc/pgx/v5/pgxpool"

    "github.com/mzzsml/minerva/internal/parser"
)

func GetHostDetails(addr string) (parser.Host, error) {
    pool := ConnectToDb()
    defer pool.Close()

    rows, err := pool.Query(context.Background(), "SELECT port.port_num, port.state, port.reason, port.service, port.version, port.extrainfo FROM port WHERE host_id = (SELECT id FROM host WHERE address = $1)", addr)
    if err != nil {
        log.Printf("ERROR: GetHostDetails: %s", err)
    }

    var p parser.Port
    var ports []parser.Port
    var h parser.Host
    for rows.Next() {
        err := rows.Scan(&p.PortId, &p.State, &p.Reason, &p.Service, &p.Version, &p.Extrainfo)
        if err != nil {
            log.Fatal(err)
        }
        ports = append(ports, p)
        //log.Printf("port %d, state %s, service %s\n", port, state, service)
    }
    h.Ports = ports
    h.Addr = addr
    return h, err
}
