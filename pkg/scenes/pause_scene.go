package scenes

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	systeminterfaces "github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/utils"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type PauseScene struct {
	sceneManager *SceneManager
	ecsManager   *ecs.ECSManager
	sceneData    *SceneData

	entities      []*ecs.Entity
	logicSystems  []systeminterfaces.UpdatableSystemInterface
	renderSystems []systeminterfaces.RenderableSystemInterface
}

func NewPauseScene(sm *SceneManager, em *ecs.ECSManager, sd *SceneData) *PauseScene {
	return &PauseScene{
		sceneManager:  sm,
		ecsManager:    em,
		sceneData:     sd,
		entities:      []*ecs.Entity{},
		logicSystems:  []systeminterfaces.UpdatableSystemInterface{},
		renderSystems: []systeminterfaces.RenderableSystemInterface{},
	}
}

// TODO: Immplement Init function
func (ps *PauseScene) Initialize() {
	// Initialize the scene here
}

func (ps *PauseScene) Update(dt float64) {
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

	// Cleanup if necessary
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

func (ps *PauseScene) AddEntity(eID *ecs.Entity) {
	ps.entities = append(ps.entities, eID)
}

func (ps *PauseScene) RemoveAllEntities() {
	for _, entity := range ps.entities {
		ps.ecsManager.DestroyEntity(entity.ID)
	}
	ps.entities = []*ecs.Entity{}
}
