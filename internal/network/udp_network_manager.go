// internal/network/udp_network_manager.go

package network

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type UDPNetworkManager struct {
	conn *net.UDPConn
	quit chan bool
}

const (
	MAX_PACKET_SIZE = 1024
)

func NewUDPNetworkManager(port int) (*UDPNetworkManager, error) {

	// Resolve the address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	// Create the connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	return &UDPNetworkManager{
		conn: conn,
		quit: make(chan bool),
	}, nil
}

func (u *UDPNetworkManager) Start() error {
	// Placeholder for any startup routines
	return nil
}

func (u *UDPNetworkManager) Send(packet *Packet, addr *net.UDPAddr) error {

	// Marshal the packet
	data, err := proto.Marshal(packet)
	if err != nil {
		return err
	}

	// Send the packet
	_, err = u.conn.WriteToUDP(data, addr)
	if err != nil {
		return err
	}

	return nil
}

func (u *UDPNetworkManager) Receive() (*Packet, *net.UDPAddr, error) {

	// Create a buffer to read the packet into
	buffer := make([]byte, MAX_PACKET_SIZE)

	// Set a deadline for reading
	u.conn.SetReadDeadline(time.Now().Add(time.Second)) // 1-Second timemout

	// Read the packet
	bytes, addr, err := u.conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, nil, err
	}

	// Unmarshal the packet
	packet := &Packet{}

	err = proto.Unmarshal(buffer[:bytes], packet)
	if err != nil {
		return nil, nil, err
	}

	return packet, addr, nil

}

func (u *UDPNetworkManager) Stop() error {
	close(u.quit)
	return u.conn.Close()
}
