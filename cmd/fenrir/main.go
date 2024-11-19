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
	sceneManager.PushScene("assets/scenes/main_menu.json")

	// Disable the Escape key from closing the window
	raylib.SetExitKey(0)

	// Main game loop
	for !raylib.WindowShouldClose() && !sceneManager.ShouldExitGame() {

		// Get frame time
		deltaTime := raylib.GetFrameTime()

		currentScene := sceneManager.GetCurrentScene()

		if currentScene != nil {

			// Update Current Scene
			currentScene.Update(float64(deltaTime))

			ecsManager.UpdateLogicSystems(float64(deltaTime))
		}

		// Apply any pending scene changes
		if sceneManager.ShouldChangeScene() {
			err := sceneManager.ApplyPendingSceneChange()
			if err != nil {
				utils.ErrorLogger.Fatalf("Failed to apply pending scene change: %v", err)
			}
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		if currentScene != nil {
			currentScene.Render()
		}

		// Render Current Scene
		sceneManager.Render()

		raylib.EndDrawing()

	}
}
