package main

import (
	"net"
	"time"

	"github.com/webbelito/Fenrir/internal/network"
	"github.com/webbelito/Fenrir/pkg/utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize the logger
	utils.InfoLogger.Println("Starting Fenrir Client")

	// Initialize the network manager

	//  Resolve the host address
	// TODO: Implement a way to specify the server address
	serverAddr, err := net.ResolveUDPAddr("udp", "localhost:8000")
	if err != nil {
		utils.FatalLogger.Fatalf("Failed to resolve server address: %s", err)
	}

	// Dial to the server
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		utils.FatalLogger.Fatalf("Failed to dial to server: %s", err)
	}
	defer conn.Close()

	// Create UDPNetworkManager
	udpNetworkManager := &network.UDPNetworkManager{
		Conn: conn,
		Quit: make(chan bool),
	}

	// Start the network manager
	err = udpNetworkManager.Start()
	if err != nil {
		utils.FatalLogger.Fatalf("Failed to start network manager: %s", err)
	}
	defer udpNetworkManager.Stop()

	// Send JoinRequest
	joinReq := &network.JoinRequest{
		PlayerName: "Webbelito",
	}

	// Create JoinRequest packet
	joinPacket := &network.Packet{
		Type: network.PacketType_JOIN,
		Payload: &network.Packet_JoinRequest{
			JoinRequest: joinReq,
		},
	}

	// Send the JoinRequest packet
	err = udpNetworkManager.Send(joinPacket)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to send JOIN packet: %s", err)
	}

	utils.InfoLogger.Println("JoinRequest sent to server")

	// Define variables to store the WelcomeResponse
	var welcomeMessage string

	// Set a timeout for receiving the WelcomeResponse
	timeout := time.After(5 * time.Second)
	done := make(chan bool)

	// Listen for incoming packets
	go func() {
		for {
			packet, addr, err := udpNetworkManager.Receive()
			if err != nil {
				// TODO: Handle timeout or other errors (continue listening)
				continue
			}

			// Check if the packet is from the server
			if addr.String() != serverAddr.String() {
				utils.WarnLogger.Printf("Received packet from unexpected address: %s", addr.String())
				continue
			}

			// Process the packet based on the packet type
			switch packet.Type {
			case network.PacketType_WELCOME:
				welcomeResp := packet.GetWelcomeResponse()
				if welcomeResp == nil {
					utils.WarnLogger.Println("Received WELCOME packet with invalid payload")
					continue
				}

				welcomeMessage = welcomeResp.Message

				utils.InfoLogger.Printf("Received WelcomeResponse: %s", welcomeMessage)

				done <- true

				// TODO: Temporary solution to break the loop
				return

			default:
				utils.WarnLogger.Printf("Received packet with unknown type: %v", packet.Type)
			}
		}
	}()

	select {
	case <-done:
		// TODO: Proceed to display the welcome message
	case <-timeout:
		utils.ErrorLogger.Println("Timed out waiting for WelcomeResponse")
	}

	// Initialize Raylib
	utils.InfoLogger.Println("Initializing Raylib")

	rl.InitWindow(800, 600, "Fenrir Client")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// Display the welcome message
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawText(welcomeMessage, 200, 300, 20, rl.Maroon)

		rl.EndDrawing()
	}
}
