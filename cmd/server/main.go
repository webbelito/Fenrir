package main

import (
	"fmt"

	"github.com/webbelito/Fenrir/internal/network"
	"github.com/webbelito/Fenrir/pkg/utils"
)

func main() {

	// Initialize the logger
	utils.InfoLogger.Println("Starting Fenrir Server")

	// Initialize the network manager
	networkManager, err := network.NewUDPNetworkManager(8000)
	if err != nil {
		utils.FatalLogger.Fatalf("Failed to create network manager: %s", err)
	}

	// Start the network manager
	err = networkManager.Start()
	if err != nil {
		utils.FatalLogger.Fatalf("Failed to start network manager: %s", err)
	}
	defer networkManager.Stop()

	// Listen for incoming packets
	for {
		packet, addr, err := networkManager.Receive()
		if err != nil {
			// TODO: Handle timeout or other errors (continue listening)
			continue
		}

		// Process the packet based on the packet type
		switch packet.Type {
		case network.PacketType_JOIN:
			joinReq := packet.GetJoinRequest()
			if joinReq == nil {
				utils.WarnLogger.Println("Received JOIN packet with invalid payload")
				continue
			}

			utils.InfoLogger.Printf("Player %s joined the game from %s", joinReq.PlayerName, addr.String())

			// Create WelcomeResponse
			// TODO: Implement a Unique PlayerID
			welcomeResp := &network.WelcomeResponse{
				PlayerId: 1,
				Message:  fmt.Sprintf("Welcome to Fenrir, %s!", joinReq.PlayerName),
			}

			// Create WelcomeResponse packet
			welcomePacket := &network.Packet{
				Type: network.PacketType_WELCOME,
				Payload: &network.Packet_WelcomeResponse{
					WelcomeResponse: welcomeResp,
				},
			}

			// Send the WelcomeResponse packet
			err = networkManager.SendTo(welcomePacket, addr)
			if err != nil {
				utils.ErrorLogger.Printf("Failed to send WelcomeResponse to %s: %s", addr.String(), err)
				continue
			}

			utils.InfoLogger.Printf("Sent WelcomeResposne to Player ID %d at %s", welcomeResp.PlayerId, addr.String())

		default:
			utils.WarnLogger.Printf("Received packet with unknown type: %v", packet.Type)
		}
	}
}
