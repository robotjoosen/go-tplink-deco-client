package model

import "net"

type (
	Clients []Client
	Client  struct{}
)

type (
	Devices []Device
	Device  struct {
		ID         string           `json:"id"`
		Name       string           `json:"name"`
		IPAddress  net.IP           `json:"ip_address"`
		MACAddress net.HardwareAddr `json:"mac_address"`
	}
)
