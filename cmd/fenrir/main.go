package main

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/ecs/systems"
	"github.com/webbelito/Fenrir/pkg/editor"
	"github.com/webbelito/Fenrir/pkg/graphics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize Raylib
	rl.InitWindow(1920, 1280, "Fenrir Test")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	renderSystem := graphics.NewRenderSystem(rl.Rectangle{X: 0, Y: 0, Width: float32(rl.GetScreenWidth()), Height: float32(rl.GetScreenHeight())})

	// Initialize ECS Manager
	ecsManager := ecs.NewECSManager()
	ecsManager.AddSystem(&systems.InputSystem{}, 0)
	ecsManager.AddSystem(&systems.MovementSystem{}, 1)
	ecsManager.AddSystem(renderSystem, 2)

	// Create a player entity
	player := ecsManager.CreateEntity()
	ecsManager.AddComponent(player, ecs.PositionComponent, &ecs.Position{Vector: rl.NewVector2(100, 100)})
	ecsManager.AddComponent(player, ecs.VelocityComponent, &ecs.Velocity{Vector: rl.NewVector2(0, 0)})
	ecsManager.AddComponent(player, ecs.ColorComponent, &ecs.Color{Color: rl.Red})
	ecsManager.AddComponent(player, ecs.SpeedComponent, &ecs.Speed{Value: 200})
	ecsManager.AddComponent(player, ecs.PlayerComponent, &ecs.Player{Name: "Webbelito"})

	// Initialize the Editor
	gameEditor := editor.NewEditor(ecsManager)

	// Main game loop

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		// Handle editor input
		handleEditorInput(gameEditor)

		// Handle spawner input
		handleSpawnerInput(ecsManager)

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
		gameEditor.Draw(&graphics.PerformanceMonitorData{FPS: rl.GetFPS(), UpdateDuration: updateDuration, RenderDuration: renderDuration, TotalDuration: totalDuration})

		// End drawing
		rl.EndDrawing()

	}
}

func handleEditorInput(editor *editor.Editor) {
	if rl.IsKeyPressed(rl.KeyF1) {
		editor.ToggleVisibility()
	}
}

func handleSpawnerInput(ecsManager *ecs.ECSManager) {

	// Create a slice of raylib colors
	colors := []rl.Color{
		rl.Blue,
		rl.Green,
		rl.Purple,
		rl.Orange,
		rl.Pink,
		rl.Yellow,
		rl.SkyBlue,
		rl.Lime,
		rl.Gold,
		rl.Violet,
		rl.Brown,
		rl.LightGray,
		rl.DarkGray,
	}

	if rl.IsKeyPressed(rl.KeySpace) {

		// SPAWN 100 ENTITIES
		// Select a random color from the colors slice

		for i := 0; i < 450; i++ {

			color := colors[rl.GetRandomValue(0, int32(len(colors)-1))]

			entity := ecsManager.CreateEntity()
			ecsManager.AddComponent(entity, ecs.PositionComponent, &ecs.Position{Vector: rl.NewVector2(float32(rl.GetRandomValue(0, int32(rl.GetScreenWidth())-1)), float32(rl.GetRandomValue(0, int32(rl.GetScreenHeight())-1)))})
			ecsManager.AddComponent(entity, ecs.VelocityComponent, &ecs.Velocity{Vector: rl.NewVector2(float32(rl.GetRandomValue(-10, 10)), float32(rl.GetRandomValue(-10, 10)))})
			ecsManager.AddComponent(entity, ecs.SpeedComponent, &ecs.Speed{Value: float32(rl.GetRandomValue(50, 200))})
			ecsManager.AddComponent(entity, ecs.ColorComponent, &ecs.Color{Color: color})
		}
	}
}
