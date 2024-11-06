// internal/network/network_manager.go

package network

import (
	"net"

	"github.com/webbelito/Fenrir/internal/network/protocol"
)

// NetworkManager defines the functionality for network commuinication
type NetworkManager interface {
	Start() error
	Send(packet *protocol.Packet, addr *net.UDPAddr) error
	Receive() (*protocol.Packet, *net.UDPAddr, error)
	Stop() error
}
