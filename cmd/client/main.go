package main

import (
	"net"
	"time"

	// ECS
	"github.com/webbelito/Fenrir/internal/ecs"
	"github.com/webbelito/Fenrir/internal/ecs/components"
	"github.com/webbelito/Fenrir/internal/ecs/systems"

	// Network
	"github.com/webbelito/Fenrir/internal/network"

	// Utils
	"github.com/webbelito/Fenrir/pkg/utils"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// TODO: Move the GameState to a seperate logic
type GameState struct {
	WelcomeMsg string
	PlayerId   uint32
	Players    map[uint32]components.Position
}

func main() {

	// Initialize the logger
	utils.InfoLogger.Println("Starting Fenrir Client")

	// Initialize GameState
	gameState := &GameState{
		Players: make(map[uint32]components.Position),
	}

	// Initialize the ECS Manager (TODO: optional on clients if client manages its own entities)
	ecsManager := ecs.NewManager()

	// Create and add ECS systems
	movementSystem := &systems.MovementSystem{}
	renderSystem := &systems.RenderSystem{}
	ecsManager.AddSystem(movementSystem)
	ecsManager.AddSystem(renderSystem)

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
	networkManager := network.NewUDPNetworkManagerWithConn(conn)

	// Start the network manager
	err = networkManager.Start()
	if err != nil {
		utils.FatalLogger.Fatalf("Failed to start network manager: %s", err)
	}
	defer networkManager.Stop()

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
	err = networkManager.Send(joinPacket)
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
			packet, addr, err := networkManager.Receive()
			if err != nil {
				// TODO: Handle timeout or other errors (continue listening)
				continue
			}

			// Handle ACK packets
			if packet.Type == network.PacketType_ACK {
				networkManager.HandleACK(packet)
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

				// Update the GameState
				gameState.WelcomeMsg = welcomeResp.Message
				gameState.PlayerId = welcomeResp.PlayerId

				utils.InfoLogger.Printf("Received WelcomeResponse: %s", welcomeMessage)

				// Initialize player's position
				ecsManager.AddEntity(ecs.Entity(gameState.PlayerId))
				ecsManager.AddComponent(ecs.Entity(gameState.PlayerId), &components.Position{X: 400, Y: 300})
				ecsManager.AddComponent(ecs.Entity(gameState.PlayerId), &components.Velocity{X: 0, Y: 0})
				ecsManager.AddComponent(ecs.Entity(gameState.PlayerId), &components.Player{Name: joinReq.PlayerName})

				// Add the player to the GameState
				gameState.Players[gameState.PlayerId] = components.Position{X: 400, Y: 300}

				// Send ACK for the WelcomeResponse
				ackPacket := &network.Packet{
					Type:        network.PacketType_ACK,
					Reliability: network.Reliability_UNRELIABLE,
					MessageId:   packet.MessageId,
					Payload: &network.Packet_Ack{
						Ack: &network.Ack{
							MessageId: packet.MessageId,
						},
					},
				}

				// Send the ACK packet
				err = networkManager.Send(ackPacket)
				if err != nil {
					utils.ErrorLogger.Printf("Failed to send ACK packet: %s", err)
				} else {
					utils.InfoLogger.Printf("Sent ACK for message ID %d", packet.MessageId)
				}

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

	// Handle incoming packets (Non-blocking)
	go func() {
		for {
			packet, addr, err := networkManager.Receive()
			if err != nil {
				// Handle timeout or other errors (continue listening)
				continue
			}

			// Handle ACK packets
			if packet.Type == network.PacketType_ACK {
				networkManager.HandleACK(packet)
				continue
			}

			// Since the client is connected, all received packets should be from the server.
			if addr.String() != serverAddr.String() {
				utils.WarnLogger.Printf("Received packet from unexpected address: %s", addr.String())
				continue
			}

			// Process the packet based on the packet type
			switch packet.Type {
			case network.PacketType_POSITION_UPDATE:
				// Handle PositionUpdate updates from server
				positionUpdate := packet.GetPositionUpdate()
				if positionUpdate == nil {
					utils.WarnLogger.Println("Received POSITION_UPDATE packet with invalid payload")
					continue
				}

				// Update game state accordingly
				gameState.Players[positionUpdate.PlayerId] = components.Position{X: positionUpdate.X, Y: positionUpdate.Y}

				// Send ACK for the PositionUpdate
				ackPacket := &network.Packet{
					Type:        network.PacketType_ACK,
					Reliability: network.Reliability_UNRELIABLE,
					MessageId:   packet.MessageId,
					Payload: &network.Packet_Ack{
						Ack: &network.Ack{
							MessageId: packet.MessageId,
						},
					},
				}

				// Send the ACK packet
				err = networkManager.Send(ackPacket)
				if err != nil {
					utils.ErrorLogger.Printf("Failed to send ACK packet: %s", err)
				} else {
					utils.InfoLogger.Printf("Sent ACK for message ID %d", packet.MessageId)
				}

			default:
				utils.WarnLogger.Printf("Received packet with unknown type: %v", packet.Type)
			}
		}
	}()

	// Game Loop
	// TODO: Game Loop should be moved to a seperate logic
	for !rl.WindowShouldClose() {

		// Calculate deltaTime
		deltaTime := rl.GetFrameTime()

		// Handle Input
		var movementX float32
		var movementY float32

		if rl.IsKeyDown(rl.KeyRight) {
			movementX += 100 * deltaTime
		}

		if rl.IsKeyDown(rl.KeyLeft) {
			movementX -= 100 * deltaTime
		}

		if rl.IsKeyDown(rl.KeyDown) {
			movementY += 100 * deltaTime
		}

		if rl.IsKeyDown(rl.KeyUp) {
			movementY -= 100 * deltaTime
		}

		// Update the Player Velocity based on input
		posComponent := ecsManager.GetComponent(ecs.Entity(gameState.PlayerId), &components.Position{}).(*components.Position)
		velComponent := ecsManager.GetComponent(ecs.Entity(gameState.PlayerId), &components.Velocity{}).(*components.Velocity)

		velComponent.X = movementX
		velComponent.Y = movementY

		// Update the ECS systems
		ecsManager.Update(deltaTime)

		// Retrieve the updated position after systems update
		updatedLocalPosition := ecsManager.GetComponent(ecs.Entity(gameState.PlayerId), &components.Position{}).(*components.Position)

		// Update local game state
		gameState.Players[gameState.PlayerId] = components.Position{X: updatedLocalPosition.X, Y: updatedLocalPosition.Y}

		// Send updated position to server as Reliable packet
		positionUpdate := &network.PositionUpdate{
			PlayerId: gameState.PlayerId,
			X:        posComponent.X,
			Y:        posComponent.Y,
		}

		positionUpdatePacket := &network.Packet{
			Type:        network.PacketType_POSITION_UPDATE,
			Reliability: network.Reliability_RELIABLE,
			Payload: &network.Packet_PositionUpdate{
				PositionUpdate: positionUpdate,
			},
		}

		err = networkManager.Send(positionUpdatePacket)
		if err != nil {
			utils.ErrorLogger.Printf("Failed to send PositionUpdate packet: %s", err)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Render Players
		// TODO, Replace with ECS RenderSystem
		/*
			for id, pos := range gameState.Players {
				rl.DrawCircle(int32(pos.X), int32(pos.Y), 25, rl.Red)
				rl.DrawText(fmt.Sprintf("ID: %d", id), int32(pos.X)-10, int32(pos.Y)-20, 24, rl.Black)
			}

			//rl.DrawText(welcomeMessage, 200, 300, 20, rl.Maroon)
		*/

		rl.EndDrawing()
	}
}
