package storage

import (
    "context"
    "log"

    "github.com/mzzsml/minerva/internal/model"
)

func (db *Db) GetHostDetails(addr string) (model.Host, error) {
    q := `
    SELECT 
        port.port_num,
        port.state,
        port.reason,
        port.service,
        port.version,
        port.extrainfo
    FROM port
    WHERE host_id = (SELECT id FROM host WHERE address = $1)`
    rows, err := db.Pool.Query(context.Background(), q, addr)
    if err != nil {
        log.Printf("ERROR: GetHostDetails: %s", err)
    }
    defer rows.Close()

    var p model.Port
    var ports []model.Port
    var h model.Host
    for rows.Next() {
        err := rows.Scan(&p.PortId, &p.State, &p.Reason, &p.Service, &p.Version, &p.Extrainfo)
        if err != nil {
            log.Fatal(err)
        }
        ports = append(ports, p)
    }
    h.Ports = ports
    h.Addr = addr
    return h, err
}
