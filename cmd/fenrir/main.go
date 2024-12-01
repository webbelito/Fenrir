package main

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/scenes"
	"github.com/webbelito/Fenrir/pkg/systems"
	"github.com/webbelito/Fenrir/pkg/utils"

	"github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize raylib
	raylib.InitWindow(1920, 1280, "Fenrir Engine")
	defer raylib.CloseWindow()

	// Initialize audio device
	raylib.InitAudioDevice()
	defer raylib.CloseAudioDevice()

	raylib.SetTargetFPS(60)

	// Initialize raygui
	raygui.LoadStyleDefault()

	raygui.SetStyle(raylib.FontDefault, raygui.TEXT_SIZE, 20)

	// Initialize ECS Manager
	ecsManager := ecs.NewManager()

	// Camera Component
	cameraEntity := ecsManager.CreateEntity()
	ecsManager.AddComponent(cameraEntity.ID, ecs.CameraComponent, &components.Camera{
		OwnerEntity: cameraEntity.ID,
		Target:      raylib.Vector2{X: 0, Y: 0},
		Offset:      raylib.Vector2{X: float32(raylib.GetScreenWidth()) / 2, Y: float32(raylib.GetScreenHeight()) / 2},
		Zoom:        1.0,
	})

	// Initialize Camera System
	cameraSystem := systems.NewCameraSystem(ecsManager, 0)
	ecsManager.RegisterSystem(cameraSystem, cameraSystem.GetPriority())

	// Initialize UI System
	UISystem := systems.NewUISystem(ecsManager, 1)
	ecsManager.RegisterSystem(UISystem, UISystem.GetPriority())

	// Initialize Event Listener System
	EventsListenerSystem := systems.NewEventsListenerSystem(ecsManager, 2)
	ecsManager.RegisterSystem(EventsListenerSystem, EventsListenerSystem.GetPriority())

	// Initialize Scene Manager
	sceneManager := scenes.NewSceneManager(ecsManager)
	err := sceneManager.PushScene("assets/scenes/main_menu.json")
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to push scene: %v", err)
	}

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

		ecsManager.Render()

		raylib.EndDrawing()

	}

	// TODO: Unload All Resources
}
