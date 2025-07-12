package model

import (
    "encoding/json"
)

type Port struct {
    Protocol  string `json:"protocol"`
    PortId    int    `json:"port"`
    State     string `json:"state"`
    Reason    string `json:"reason"`
    Product   string `json:"product"`
    Service   string `json:"service"`
    Version   string `json:"version"`
    Extrainfo string `json:"extrainfo"`
}

type Host struct {
    Addr  string `json:"addr"`
    Ports []Port `json:"ports"`
}

func (h *Host) ParseFromJson(b []byte) error {
    return json.Unmarshal(b, &h)
}
