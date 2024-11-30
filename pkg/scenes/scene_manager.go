package scenes

import (
	"fmt"

	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/events"
	"github.com/webbelito/Fenrir/pkg/utils"
)

// SceneManager is a struct that manages scenes
type SceneManager struct {
	scenes           []Scene
	manager          *ecs.Manager
	shouldExitGame   bool
	pendingScenePath string
	pendingChange    bool
}

// NewSceneManager creates a new SceneManager
func NewSceneManager(m *ecs.Manager) *SceneManager {
	sm := &SceneManager{
		scenes:         []Scene{},
		manager:        m,
		shouldExitGame: false,
	}

	// Subscribe to events
	sm.manager.SubscribeEvent("change_scene", sm.OnChangeScene)
	sm.manager.SubscribeEvent("exit_game", sm.OnExitGame)

	return sm

}

/*
PushScene pushes a new scene onto the scene stack
*/
func (sm *SceneManager) PushScene(sceneFilePath string) error {

	// Load the scene data
	sceneData, err := LoadSceneData(sceneFilePath)
	if err != nil {
		utils.ErrorLogger.Println("Failed to load scene data: ", err)
		return err
	}

	// Declare a new scene
	var newScene Scene

	// Create the scene based on the scene name
	switch sceneData.SceneName {
	case "MainMenu":
		newScene = NewMainMenuScene(sm, sm.manager, sceneData)
	case "Game":
		newScene = NewGameScene(sm, sm.manager, sceneData)
	case "Pause":
		newScene = NewPauseScene(sm, sm.manager, sceneData)
	default:
		return fmt.Errorf("scene manager: unkown scene %s", sceneData.SceneName)
	}

	// If there are scenes on the stack, pause the current scene
	if len(sm.scenes) > 0 {
		// Optional: Pause the current scene
		sm.scenes[len(sm.scenes)-1].Pause()
	}

	// Initialize the new scene
	sm.scenes = append(sm.scenes, newScene)
	newScene.Initialize()

	utils.InfoLogger.Printf("Pushed scene: %s\n", sceneData.SceneName)

	return nil
}

/*
PopScene pops the topmost scene from the scene stack
*/
func (sm *SceneManager) PopScene() error {

	// Check if there are scenes to pop
	if len(sm.scenes) == 0 {
		return fmt.Errorf("scene manager: no scenes to pop")
	}

	// Cleanup the current scene
	currentScene := sm.scenes[len(sm.scenes)-1]
	currentScene.Cleanup()
	sm.scenes = sm.scenes[:len(sm.scenes)-1]

	// Optionally resume the previous scene
	if len(sm.scenes) > 0 {
		sm.scenes[len(sm.scenes)-1].Resume()
	}

	return nil

}

/*
ChangeScene changes the current scene to a new scene
*/
func (sm *SceneManager) ChangeScene(sceneFilePath string) error {

	// Cleanup the current scene
	for _, scene := range sm.scenes {
		scene.Cleanup()
	}

	// Clear the scene stack
	sm.scenes = []Scene{}

	// Load and initialize the new scene
	return sm.PushScene(sceneFilePath)
}

/*
GetCurrentScene returns the current scene
*/
func (sm *SceneManager) GetCurrentScene() Scene {

	// Check if there are scenes
	if len(sm.scenes) == 0 {
		return nil
	}

	// Return the topmost scene
	return sm.scenes[len(sm.scenes)-1]
}

/*
SetCurrentScene sets the current scene
*/
func (sm *SceneManager) SetCurrentScene(sceneFilePath string) error {

	// Set the pending scene change
	sm.pendingScenePath = sceneFilePath
	sm.pendingChange = true

	return nil
}

// ShouldExitGame returns true if the game should exit
func (sm *SceneManager) ShouldExitGame() bool {
	return sm.shouldExitGame
}

/*
ShouldChangeScene returns true if a scene change is pending
*/
func (sm *SceneManager) ShouldChangeScene() bool {
	return sm.pendingChange
}

/*
ApplyPendingSceneChange applies the pending scene change
*/
func (sm *SceneManager) ApplyPendingSceneChange() error {

	// Check if there is a pending scene change
	if sm.pendingChange {

		// Change the scene
		err := sm.ChangeScene(sm.pendingScenePath)
		if err != nil {
			return err
		}

		// Reset the pending scene change
		sm.pendingChange = false
		sm.pendingScenePath = ""

	}

	return nil
}

func (sm *SceneManager) Update(dt float64) {

	// Check if there are scenes
	if len(sm.scenes) == 0 {
		return
	}

	// Update the topmost scene
	sm.scenes[len(sm.scenes)-1].Update(dt)
}

func (sm *SceneManager) Render() {

	// Check if there are scenes
	if len(sm.scenes) == 0 {
		return
	}

	// Render all scenes
	for _, scene := range sm.scenes {
		scene.Render()
	}
}

/*
OnChangeScene is an event handler for changing scenes
*/
func (sm *SceneManager) OnChangeScene(event events.Event) {

	// Check the type of the event
	switch e := event.(type) {

	case events.SceneChangeEvent:
		// Change the scene
		err := sm.ChangeScene(e.ScenePath)
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}

	default:
		utils.ErrorLogger.Println("SceneManager: OnChangeScene: unknown event type")
	}
}

/*
OnExitGame is an event handler for exiting the game
*/
func (sm *SceneManager) OnExitGame(event events.Event) {

	// Set the should exit game flag
	sm.shouldExitGame = event.(events.ExitGameEvent).ShouldExitGame
}
