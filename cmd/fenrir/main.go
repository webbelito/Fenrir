package main

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/graphics"

	rlgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize Raylib
	rl.InitWindow(1280, 720, "Fenrir Test")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// Initialize ECS Manager
	ecsManager := ecs.NewECSManager()
	ecsManager.AddSystem(&ecs.MovementSystem{}, 1)
	ecsManager.AddSystem(&graphics.RenderSystem{}, 2)

	// Create a player entity
	player := ecsManager.CreateEntity()
	ecsManager.AddComponent(player, ecs.PositionComponent, &ecs.Position{Vector: rl.NewVector2(100, 100)})
	ecsManager.AddComponent(player, ecs.VelocityComponent, &ecs.Velocity{Vector: rl.NewVector2(1, 1)})

	// Main game loop

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		// Handle input
		handleInput(ecsManager, player)

		// Update ECS systems
		ecsManager.Update(float64(deltaTime))

		// Begin drawing
		rl.BeginDrawing()

		rlgui.Panel(rl.NewRectangle(0, 0, 200, 720), "Controls")
		rlgui.Label(rl.NewRectangle(10, 20, 200, 20), "WASD to move")

		// Clear screen
		rl.ClearBackground(rl.RayWhite)

		// End drawing
		rl.EndDrawing()

	}
}

// TODO: Refactor to a InputSystem
func handleInput(ecsManager *ecs.ECSManager, player ecs.Entity) {
	if rl.IsKeyDown(rl.KeyW) {
		velocity := ecsManager.GetComponent(player, ecs.VelocityComponent).(*ecs.Velocity)
		velocity.Vector.Y -= 1
	}
	if rl.IsKeyDown(rl.KeyS) {
		velocity := ecsManager.GetComponent(player, ecs.VelocityComponent).(*ecs.Velocity)
		velocity.Vector.Y += 1
	}
	if rl.IsKeyDown(rl.KeyA) {
		velocity := ecsManager.GetComponent(player, ecs.VelocityComponent).(*ecs.Velocity)
		velocity.Vector.X -= 1
	}
	if rl.IsKeyDown(rl.KeyD) {
		velocity := ecsManager.GetComponent(player, ecs.VelocityComponent).(*ecs.Velocity)
		velocity.Vector.X += 1
	}
}
