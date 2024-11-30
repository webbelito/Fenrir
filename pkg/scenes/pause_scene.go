package scenes

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// PauseScene is a scene for the pause menu
type PauseScene struct {
	sceneManager *SceneManager
	manager      *ecs.Manager
	sceneData    *SceneData

	sceneEntities []*ecs.Entity
}

// NewPauseScene creates a new PauseScene
func NewPauseScene(sm *SceneManager, m *ecs.Manager, sd *SceneData) *PauseScene {
	return &PauseScene{
		sceneManager:  sm,
		manager:       m,
		sceneData:     sd,
		sceneEntities: []*ecs.Entity{},
	}
}

// TODO: Immplement Init function
func (ps *PauseScene) Initialize() {
	// Initialize the scene here
}

func (ps *PauseScene) Update(dt float64) {

	// Check if the ECS Manager is nil
	if ps.manager == nil {
		utils.ErrorLogger.Println("PauseScene: ECS Manager is nil")
		return
	}

	// Handle Pause Input
	if raylib.IsKeyPressed(raylib.KeyP) || raylib.IsKeyPressed(raylib.KeyEscape) {

		// Resume the game
		err := ps.sceneManager.PopScene()
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}
}

func (ps *PauseScene) Render() {

	// Draw semi-transparent overlay
	raylib.DrawRectangle(0, 0, int32(raylib.GetScreenWidth()), int32(raylib.GetScreenHeight()), raylib.Fade(raylib.Black, 0.5))

	// Draw Pause Menu UI
	raygui.Label(raylib.Rectangle{
		X:      float32(raylib.GetScreenWidth()/2 - 100),
		Y:      float32(raylib.GetScreenHeight()/2 - 25),
		Width:  200,
		Height: 50,
	}, "Paused")

	// Draw Resume Button
	if raygui.Button(raylib.Rectangle{
		X:      float32(raylib.GetScreenWidth()/2 - 100),
		Y:      float32(raylib.GetScreenHeight()/2 + 25),
		Width:  200,
		Height: 50,
	}, "Resume") {
		err := ps.sceneManager.PopScene()
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}

	// Draw Exit Button
	if raygui.Button(raylib.Rectangle{
		X:      float32(raylib.GetScreenWidth()/2 - 100),
		Y:      float32(raylib.GetScreenHeight()/2 + 100),
		Width:  200,
		Height: 50,
	}, "Exit Game") {
		err := ps.sceneManager.ChangeScene("assets/scenes/main_menu.json")
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}
}

// TODO: Implement Cleanup function
func (ps *PauseScene) Cleanup() {
	// Remove all entities created by this scene
	ps.RemoveAllEntities()
}

func (ps *PauseScene) Pause() {
	// TODO: Implement Pause functionality
	// Pause game logic if necessary
	// For example, stop certain systems or timers
}

func (ps *PauseScene) Resume() {
	// TODO: Implement Resume functionality
	// Resume game logic if necessary
	// For example, resume certain systems or timers
}

/*
AddEntity adds an entity to the scene entity list
Usage is to identify which entities were created by the scene
*/
func (ps *PauseScene) AddEntity(eID *ecs.Entity) {
	ps.sceneEntities = append(ps.sceneEntities, eID)
}

/*
RemoveAllEntities removes all entities created by the scene
*/
func (ps *PauseScene) RemoveAllEntities() {

	// Iterate over all entities created by the scene
	for _, entity := range ps.sceneEntities {
		ps.manager.DestroyEntity(entity.ID)
	}

	// Clear the scene entities list
	ps.sceneEntities = []*ecs.Entity{}
}
