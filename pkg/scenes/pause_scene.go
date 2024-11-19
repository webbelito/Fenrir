package scenes

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type PauseScene struct {
	sceneManager *SceneManager
	ecsManager   *ecs.ECSManager
	sceneData    *SceneData

	entities []*ecs.Entity
}

func NewPauseScene(sm *SceneManager, em *ecs.ECSManager, sd *SceneData) *PauseScene {
	return &PauseScene{
		sceneManager: sm,
		ecsManager:   em,
		sceneData:    sd,
		entities:     []*ecs.Entity{},
	}
}

// TODO: Immplement Init function
func (ps *PauseScene) Init() {
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
	raygui.Label(raylib.NewRectangle(float32(raylib.GetScreenWidth()/2-100), float32(raylib.GetScreenHeight()/2-100), 200, 50), "PAUSED")

	// Draw Resume Button
	if raygui.Button(raylib.NewRectangle(
		float32(raylib.GetScreenWidth()/2-100),
		float32(raylib.GetScreenHeight()/2),
		200,
		50,
	), "Resume") {
		err := ps.sceneManager.PopScene()
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}

	// Draw Exit Button
	if raygui.Button(raylib.NewRectangle(
		float32(raylib.GetScreenWidth()/2-100),
		float32(raylib.GetScreenHeight()/2+100),
		200,
		50,
	), "Exit to Main Menu") {
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
