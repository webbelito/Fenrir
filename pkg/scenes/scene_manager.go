package scenes

import (
	"fmt"

	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/events"
	"github.com/webbelito/Fenrir/pkg/utils"
)

type SceneManager struct {
	scenes           []Scene
	ecsManager       *ecs.ECSManager
	shouldExitGame   bool
	pendingScenePath string
	pendingChange    bool
}

func NewSceneManager(ecsManager *ecs.ECSManager) *SceneManager {
	sm := &SceneManager{
		scenes:         []Scene{},
		ecsManager:     ecsManager,
		shouldExitGame: false,
	}

	sm.ecsManager.GetEventsManager().Subscribe("change_scene", sm.OnChangeScene)
	sm.ecsManager.GetEventsManager().Subscribe("exit_game", sm.OnExitGame)

	return sm

}

func (sm *SceneManager) PushScene(sceneFilePath string) error {

	sceneData, err := LoadSceneData(sceneFilePath)
	if err != nil {
		utils.ErrorLogger.Println("Failed to load scene data: ", err)
		return err
	}

	var newScene Scene
	switch sceneData.SceneName {
	case "MainMenu":
		newScene = NewMainMenuScene(sm, sm.ecsManager, sceneData)
	case "Game":
		newScene = NewGameScene(sm, sm.ecsManager, sceneData)
	case "Pause":
		newScene = NewPauseScene(sm, sm.ecsManager, sceneData)
	default:
		return fmt.Errorf("scene manager: unkown scene %s", sceneData.SceneName)
	}

	if len(sm.scenes) > 0 {
		// Optional: Pause the current scene
		sm.scenes[len(sm.scenes)-1].Pause()
	}

	sm.scenes = append(sm.scenes, newScene)
	newScene.Initialize()

	utils.InfoLogger.Printf("Pushed scene: %s\n", sceneData.SceneName)

	return nil
}

func (sm *SceneManager) PopScene() error {
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

func (sm *SceneManager) ChangeScene(sceneFilePath string) error {

	for _, scene := range sm.scenes {
		scene.Cleanup()
	}

	sm.scenes = []Scene{}

	// Load and initialize the new scene
	return sm.PushScene(sceneFilePath)
}

func (sm *SceneManager) GetCurrentScene() Scene {
	if len(sm.scenes) == 0 {
		return nil
	}

	return sm.scenes[len(sm.scenes)-1]
}

func (sm *SceneManager) SetCurrentScene(sceneFilePath string) error {
	sm.pendingScenePath = sceneFilePath
	sm.pendingChange = true
	return nil
}

func (sm *SceneManager) ShouldExitGame() bool {
	return sm.shouldExitGame
}

func (sm *SceneManager) ExitGame() {
	sm.shouldExitGame = true
}

func (sm *SceneManager) ShouldChangeScene() bool {
	return sm.pendingChange
}

func (sm *SceneManager) ApplyPendingSceneChange() error {
	if sm.pendingChange {

		err := sm.ChangeScene(sm.pendingScenePath)
		if err != nil {
			return err
		}

		sm.pendingChange = false
		sm.pendingScenePath = ""

	}
	return nil
}

func (sm *SceneManager) Update(dt float64) {
	if len(sm.scenes) == 0 {
		return
	}

	// Update the topmost scene
	sm.scenes[len(sm.scenes)-1].Update(dt)
}

func (sm *SceneManager) Render() {
	for _, scene := range sm.scenes {
		scene.Render()
	}
}

func (sm *SceneManager) OnChangeScene(event events.Event) {

	utils.InfoLogger.Println("SceneManager: OnChangeScene: event received")

	switch e := event.(type) {
	case events.SceneChangeEvent:

		utils.InfoLogger.Printf("SceneManager: OnChangeScene: changing scene to %s\n", e.ScenePath)

		err := sm.ChangeScene(e.ScenePath)
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}

	default:
		utils.ErrorLogger.Println("SceneManager: OnChangeScene: unknown event type")
	}
}

func (sm *SceneManager) OnExitGame(event events.Event) {

	sm.shouldExitGame = event.(events.ExitGameEvent).ShouldExitGame

}
