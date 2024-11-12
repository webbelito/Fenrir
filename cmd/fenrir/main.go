package main

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/editor"
	"github.com/webbelito/Fenrir/pkg/graphics"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize Raylib
	rl.InitWindow(1920, 1280, "Fenrir Test")
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

	// Initialize the Editor
	gameEditor := editor.NewEditor(ecsManager)

	// Main game loop

	for !rl.WindowShouldClose() {
		deltaTime := rl.GetFrameTime()

		// Handle editor input
		handleEditorInput(gameEditor)

		// Handle input
		handleGameInput(ecsManager, player)

		// Handle spawner input
		handleSpawnerInput(ecsManager)

		// Update ECS systems
		ecsManager.Update(float64(deltaTime))

		// Begin drawing
		rl.BeginDrawing()

		// Render Editor Overlay
		gameEditor.Draw()

		// Clear screen
		rl.ClearBackground(rl.Black)

		// End drawing
		rl.EndDrawing()

	}
}

// TODO: Refactor to a InputSystem
func handleGameInput(ecsManager *ecs.ECSManager, player ecs.Entity) {
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

func handleEditorInput(editor *editor.Editor) {
	if rl.IsKeyPressed(rl.KeyF1) {
		editor.ToggleVisibility()
	}
}

func handleSpawnerInput(ecsManager *ecs.ECSManager) {

	if rl.IsKeyPressed(rl.KeySpace) {

		for i := 0; i < 100; i++ {

			entity := ecsManager.CreateEntity()
			ecsManager.AddComponent(entity, ecs.PositionComponent, &ecs.Position{Vector: rl.NewVector2(float32(rl.GetRandomValue(0, 1920)), float32(rl.GetRandomValue(0, 1280)))})
			ecsManager.AddComponent(entity, ecs.VelocityComponent, &ecs.Velocity{Vector: rl.NewVector2(0, 0)})

		}
	}
}
