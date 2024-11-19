package main

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/scenes"
	"github.com/webbelito/Fenrir/pkg/utils"

	"github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize raylib
	raylib.InitWindow(1920, 1280, "Fenrir Engine")
	defer raylib.CloseWindow()

	raylib.SetTargetFPS(60)

	// Initialize raygui
	raygui.LoadStyleDefault()

	// Initialize Resources Manager
	// TODO: Implement resource manager
	//resourcesManager := resources.NewResourceManager()

	// Initialize ECS Manager
	ecsManager := ecs.NewECSManager()

	// Initialize Scene Manager
	sceneManager := scenes.NewSceneManager(ecsManager)

	// Set Initial Scene to Main Menu
	err := sceneManager.ChangeScene("assets/scenes/main_menu.json")
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to change scene: %v", err)
	}

	// Disable the Escape key from closing the window
	raylib.SetExitKey(0)

	// Main game loop
	for !raylib.WindowShouldClose() {

		// Get frame time
		deltaTime := raylib.GetFrameTime()

		// Update Scene Manager
		sceneManager.Update(float64(deltaTime))

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		// Render Current Scene
		sceneManager.Render()

		raylib.EndDrawing()

	}
}
