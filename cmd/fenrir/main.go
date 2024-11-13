package main

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/editor"
	"github.com/webbelito/Fenrir/pkg/systems"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize Raylib
	rl.InitWindow(1920, 1280, "Fenrir Test")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	renderSystem := systems.NewRenderSystem(rl.Rectangle{X: 0, Y: 0, Width: float32(rl.GetScreenWidth()), Height: float32(rl.GetScreenHeight())})

	// Initialize ECS Manager
	ecsManager := ecs.NewECSManager()

	// Initialize the Editor
	gameEditor := editor.NewEditor(ecsManager)

	// Add systems to the ECS Manager
	ecsManager.AddSystem(&systems.InputSystem{Editor: gameEditor, EcsManager: ecsManager}, 0)
	ecsManager.AddSystem(&systems.MovementSystem{}, 1)
	ecsManager.AddSystem(renderSystem, 2)

	// Create a player entity
	player := ecsManager.CreateEntity()
	ecsManager.AddComponent(player, ecs.PositionComponent, &components.Position{Vector: rl.NewVector2(100, 100)})
	ecsManager.AddComponent(player, ecs.VelocityComponent, &components.Velocity{Vector: rl.NewVector2(0, 0)})
	ecsManager.AddComponent(player, ecs.ColorComponent, &components.Color{Color: rl.Red})
	ecsManager.AddComponent(player, ecs.SpeedComponent, &components.Speed{Value: 200})
	ecsManager.AddComponent(player, ecs.PlayerComponent, &components.Player{Name: "Webbelito"})

	// Main game loop

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		// Update ECS entities
		updateStart := time.Now()
		ecsManager.Update(float64(deltaTime))
		updateDuration := time.Since(updateStart)

		// Begin drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Render ECS entities
		renderStart := time.Now()
		ecsManager.Render()
		renderDuration := time.Since(renderStart)

		// Calculate the total time taken for the update and render steps
		totalDuration := time.Since(updateStart)

		// Render Editor Overlay
		gameEditor.Draw(&editor.PerformanceMonitorData{FPS: rl.GetFPS(), UpdateDuration: updateDuration, RenderDuration: renderDuration, TotalDuration: totalDuration})

		// End drawing
		rl.EndDrawing()

	}
}
