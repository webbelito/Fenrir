// internal/network/udp_network_manager.go

package network

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/webbelito/Fenrir/pkg/utils"
	"google.golang.org/protobuf/proto"
)

type PendingMessage struct {
	packet   *Packet
	addr     *net.UDPAddr
	retries  int
	lastSent time.Time
}

type UDPNetworkManager struct {
	Conn               *net.UDPConn
	Quit               chan bool
	PendingMessages    map[uint64]PendingMessage
	PendingMutex       sync.Mutex
	MessageIDGenerator *MessageIDGenerator
	RetransmissionTTL  time.Duration
	MaxRetries         int
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
		Conn:               conn,
		Quit:               make(chan bool),
		PendingMessages:    make(map[uint64]PendingMessage),
		MessageIDGenerator: NewMessageIDGenerator(),
		RetransmissionTTL:  time.Second * 5,
		MaxRetries:         5,
	}, nil
}

func NewUDPNetworkManagerWithConn(conn *net.UDPConn) *UDPNetworkManager {
	return &UDPNetworkManager{
		Conn:               conn,
		Quit:               make(chan bool),
		PendingMessages:    make(map[uint64]PendingMessage),
		MessageIDGenerator: NewMessageIDGenerator(),
		RetransmissionTTL:  time.Second * 5,
		MaxRetries:         5,
	}
}

func (u *UDPNetworkManager) Start() error {
	go u.retransmissionLoop()

	return nil
}

// Applicable for clients
func (u *UDPNetworkManager) Send(packet *Packet) error {

	if packet.Reliability != Reliability_RELIABLE {

		// Unreliable packets are sent immediately without tracking

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

	// Assign a unique message ID
	packet.MessageId = u.MessageIDGenerator.NextID()

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

	// Add the packet to the pending messages
	u.PendingMutex.Lock()

	u.PendingMessages[packet.MessageId] = PendingMessage{
		packet:   packet,
		addr:     nil,
		retries:  0,
		lastSent: time.Now(),
	}

	u.PendingMutex.Unlock()

	return nil
}

// Applicable for servers
func (u *UDPNetworkManager) SendTo(packet *Packet, addr *net.UDPAddr) error {

	// Unreliable packets are sent immediately without tracking
	if packet.Reliability != Reliability_RELIABLE {

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

	// Assign a unique message ID
	packet.MessageId = u.MessageIDGenerator.NextID()

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

	// Add the packet to the pending messages
	u.PendingMutex.Lock()
	defer u.PendingMutex.Unlock()

	u.PendingMessages[packet.MessageId] = PendingMessage{
		packet:   packet,
		addr:     addr,
		retries:  0,
		lastSent: time.Now(),
	}

	return nil

}

func (u *UDPNetworkManager) retransmissionLoop() {

	// Create a ticker for retransmissions
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// Loop until the Quit channel is closed
	for {
		select {
		case <-u.Quit:
			return
		case <-ticker.C:
			u.handleRetransmissions()
		}
	}
}

func (u *UDPNetworkManager) handleRetransmissions() {

	now := time.Now()

	u.PendingMutex.Lock()
	defer u.PendingMutex.Unlock()

	for msgID, pending := range u.PendingMessages {
		if now.Sub(pending.lastSent) > u.RetransmissionTTL {
			if pending.retries >= u.MaxRetries {
				utils.ErrorLogger.Printf("Max retries reached for message ID: %d. Giving up.", msgID)
				delete(u.PendingMessages, msgID)
				continue
			}

			// Retransmit the packet
			var err error

			// Check if we should use Send or SendTo
			if pending.addr != nil {
				err = u.SendTo(pending.packet, pending.addr)
			} else {
				err = u.Send(pending.packet)
			}

			if err != nil {
				utils.ErrorLogger.Printf("Failed to retransmit message ID: %d: %s", msgID, err)
				continue
			}

			// Update the pending message
			pending.retries++
			pending.lastSent = now

			utils.InfoLogger.Printf("Retransmitted message ID: %d (Attempt: %d)", msgID, pending.retries)

		}
	}
}

func (u *UDPNetworkManager) HandleACK(packet *Packet) {

	if packet.Type != PacketType_ACK {
		return
	}

	ack := packet.GetAck()
	if ack == nil {
		utils.ErrorLogger.Println("Received ACK packet with invalid payload")
		return
	}

	u.PendingMutex.Lock()
	defer u.PendingMutex.Unlock()

	// Remove the message from the pending messages
	if _, exists := u.PendingMessages[ack.MessageId]; exists {
		// ACK received, remove the message
		delete(u.PendingMessages, ack.MessageId)
		utils.InfoLogger.Printf("Received ACK for message ID: %d", ack.MessageId)
	}
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
