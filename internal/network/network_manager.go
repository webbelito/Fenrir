// internal/network/network_manager.go

package network

import (
	"net"
)

// NetworkManager defines the functionality for network commuinication
type NetworkManager interface {
	Start() error
	Send(packet *Packet, addr *net.UDPAddr) error
	Receive() (*Packet, *net.UDPAddr, error)
	Stop() error
}
