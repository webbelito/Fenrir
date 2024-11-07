// internal/network/udp_network_manager.go

package network

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/protobuf/proto"
)

type UDPNetworkManager struct {
	Conn *net.UDPConn
	Quit chan bool
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
		Conn: conn,
		Quit: make(chan bool),
	}, nil
}

func (u *UDPNetworkManager) Start() error {
	// Placeholder for any startup routines
	return nil
}

// Applicable for clients
func (u *UDPNetworkManager) Send(packet *Packet) error {

	// Marshal the packet
	data, err := proto.Marshal(packet)
	if err != nil {
		return err
	}

	// Send the packet
	_, err = u.Conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// Applicable for servers
func (u *UDPNetworkManager) SendTo(packet *Packet, addr *net.UDPAddr) error {

	// Marshal the packet
	data, err := proto.Marshal(packet)
	if err != nil {
		return err
	}

	// Send the packet
	_, err = u.Conn.WriteToUDP(data, addr)
	if err != nil {
		return err
	}

	return nil
}

func (u *UDPNetworkManager) Receive() (*Packet, *net.UDPAddr, error) {

	// Create a buffer to read the packet into
	buffer := make([]byte, MAX_PACKET_SIZE)

	// Set a deadline for reading
	u.Conn.SetReadDeadline(time.Now().Add(time.Second)) // 1-Second timemout

	// Read the packet
	bytes, addr, err := u.Conn.ReadFromUDP(buffer)
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
	close(u.Quit)
	return u.Conn.Close()
}
