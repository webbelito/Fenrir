package main

import (
	"fmt"

	// ECS
	"github.com/webbelito/Fenrir/internal/ecs"
	"github.com/webbelito/Fenrir/internal/ecs/components"
	"github.com/webbelito/Fenrir/internal/ecs/systems"

	// Network
	"github.com/webbelito/Fenrir/internal/network"

	// Utils
	"github.com/webbelito/Fenrir/pkg/utils"
)

func main() {

	// Initialize the logger
	utils.InfoLogger.Println("Starting Fenrir Server")

	// Initialize the ECS ecsManager
	ecsManager := ecs.NewManager()

	// Create and add ECS systems
	movementSystem := &systems.MovementSystem{}
	ecsManager.AddSystem(movementSystem)

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

		// Handle ACK packets
		if packet.Type == network.PacketType_ACK {
			networkManager.HandleACK(packet)
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

			// Create a new player entity
			playerEntity := ecs.NewEntity()
			ecsManager.AddEntity(playerEntity)

			// Add components
			ecsManager.AddComponent(playerEntity, &components.Position{X: 0, Y: 0})
			ecsManager.AddComponent(playerEntity, &components.Velocity{X: 1, Y: 1})
			ecsManager.AddComponent(playerEntity, &components.Player{Name: joinReq.PlayerName})

			// Create WelcomeResponse
			// TODO: Implement a Unique PlayerID
			welcomeResp := &network.WelcomeResponse{
				PlayerId: 1,
				Message:  fmt.Sprintf("Welcome to Fenrir, %s!", joinReq.PlayerName),
			}

			// Create WelcomeResponse packet
			welcomePacket := &network.Packet{
				Type:        network.PacketType_WELCOME,
				Reliability: network.Reliability_RELIABLE,
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

		case network.PacketType_POSITION_UPDATE:
			updatePosition := packet.GetPositionUpdate()
			if updatePosition == nil {
				utils.WarnLogger.Println("Received POSITION_UPDATE packet with invalid payload")
				continue
			}

			// Update game state accordingly
			ecsManager.AddEntity(ecs.Entity(updatePosition.PlayerId))
			ecsManager.AddComponent(ecs.Entity(updatePosition.PlayerId), &components.Position{X: updatePosition.X, Y: updatePosition.Y})
			ecsManager.AddComponent(ecs.Entity(updatePosition.PlayerId), &components.Velocity{X: 0, Y: 0})
			ecsManager.AddComponent(ecs.Entity(updatePosition.PlayerId), &components.Player{Name: string(updatePosition.PlayerId)})

			utils.InfoLogger.Printf("Received POSITION_UPDATE for Player ID %d: (%f, %f)", updatePosition.PlayerId, updatePosition.X, updatePosition.Y)

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
			err = networkManager.SendTo(ackPacket, addr)
			if err != nil {
				utils.ErrorLogger.Printf("Failed to send ACK packet: %s", err)
			} else {
				utils.InfoLogger.Printf("Sent ACK for message ID %d", packet.MessageId)
			}

		default:
			utils.WarnLogger.Printf("Received packet with unknown type: %v", packet.Type)
		}

		// Update the ECS systems
		ecsManager.Update(1.0 / 60.0) // Assume 60 FPS
	}
}
