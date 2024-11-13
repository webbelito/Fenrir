package main

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/editor"
	"github.com/webbelito/Fenrir/pkg/systems"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize Raylib
	raylib.InitWindow(1920, 1280, "Fenrir Test")
	defer raylib.CloseWindow()

	raylib.SetTargetFPS(60)

	// Initialize ECS Manager
	ecsManager := ecs.NewECSManager()

	// Initialize the Editor
	gameEditor := editor.NewEditor(ecsManager)

	// Add systems to the ECS Manager
	ecsManager.AddSystem(&systems.InputSystem{
		Editor:     gameEditor,
		EcsManager: ecsManager,
	}, 0)

	ecsManager.AddSystem(&systems.MovementSystem{}, 1)

	// TODO: Implement a camera and follow the position of the camera with the screen culling rect
	ecsManager.AddSystem(&systems.RenderSystem{ScreenCullingRect: raylib.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(raylib.GetScreenWidth()),
		Height: float32(raylib.GetScreenHeight()),
	}}, 2)

	// Create a player entity
	player := ecsManager.CreateEntity()

	// Add components to the player entity
	ecsManager.AddComponent(player, ecs.PositionComponent, &components.Position{Vector: raylib.NewVector2(100, 100)})
	ecsManager.AddComponent(player, ecs.VelocityComponent, &components.Velocity{Vector: raylib.NewVector2(0, 0)})
	ecsManager.AddComponent(player, ecs.ColorComponent, &components.Color{Color: raylib.Red})
	ecsManager.AddComponent(player, ecs.SpeedComponent, &components.Speed{Value: 200})
	ecsManager.AddComponent(player, ecs.PlayerComponent, &components.Player{Name: "Webbelito"})

	// Main game loop
	for !raylib.WindowShouldClose() {

		// Get the time taken for the last frame
		deltaTime := raylib.GetFrameTime()

		// Update ECS entities
		updateStart := time.Now()
		ecsManager.Update(float64(deltaTime))
		updateDuration := time.Since(updateStart)

		// Begin drawing
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		// Render ECS entities
		renderStart := time.Now()
		ecsManager.Render()
		renderDuration := time.Since(renderStart)

		// Calculate the total time taken for the update and render steps
		totalDuration := time.Since(updateStart)

		// Render Editor Overlay
		gameEditor.Draw(&editor.PerformanceMonitorData{FPS: raylib.GetFPS(), UpdateDuration: updateDuration, RenderDuration: renderDuration, TotalDuration: totalDuration})

		// End drawing
		raylib.EndDrawing()

	}
}
