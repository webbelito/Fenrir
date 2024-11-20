package scenes

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/utils"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type MainMenuScene struct {
	sceneManager *SceneManager
	ecsManager   *ecs.ECSManager
	sceneData    *SceneData

	entities      []*ecs.Entity
	logicSystems  []systeminterfaces.Updatable
	renderSystems []systeminterfaces.Renderable
}

func NewMainMenuScene(sm *SceneManager, em *ecs.ECSManager, sd *SceneData) *MainMenuScene {
	return &MainMenuScene{
		sceneManager:  sm,
		ecsManager:    em,
		sceneData:     sd,
		entities:      []*ecs.Entity{},
		logicSystems:  []systeminterfaces.Updatable{},
		renderSystems: []systeminterfaces.Renderable{},
	}
}

// TODO: Immplement Init function
func (mms *MainMenuScene) Init() {
	// Initialize the scene here
}

func (mms *MainMenuScene) Update(dt float64) {
	// TODO: Remove temporary input handling
	if raylib.IsKeyPressed(raylib.KeyEnter) {
		err := mms.sceneManager.ChangeScene("assets/scenes/game_scene.json")
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}

	// TODO: Remove temporary input handling
	if raylib.IsKeyPressed(raylib.KeyEscape) {
		mms.sceneManager.ExitGame()
	}
}

func (mms *MainMenuScene) Render() {

	raygui.Label(raylib.Rectangle{
		X:      float32(raylib.GetScreenWidth()/2 - 100),
		Y:      float32(raylib.GetScreenHeight()/2 - 25),
		Width:  200,
		Height: 50,
	}, "Fenerir Engine")

	if raygui.Button(raylib.Rectangle{
		X:      float32(raylib.GetScreenWidth()/2 - 100),
		Y:      float32(raylib.GetScreenHeight()/2 + 25),
		Width:  200,
		Height: 50,
	}, "Start Game") {

		err := mms.sceneManager.ChangeScene("assets/scenes/game_scene.json")
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}

	if raygui.Button(raylib.Rectangle{
		X:      float32(raylib.GetScreenWidth()/2 - 100),
		Y:      float32(raylib.GetScreenHeight()/2 + 100),
		Width:  200,
		Height: 50,
	}, "Exit Game") {
		mms.sceneManager.ExitGame()
	}

}

// TODO: Implement Cleanup function
func (mms *MainMenuScene) Cleanup() {
	// Remove all entities created by this scene
	mms.RemoveAllEntities()

	// Cleanup if necessary
}

func (mms *MainMenuScene) Pause() {
	// TODO: Implement Pause functionality
	// Pause game logic if necessary
	// For example, stop certain systems or timers
}

func (mms *MainMenuScene) Resume() {
	// TODO: Implement Resume functionality
	// Resume game logic if necessary
	// For example, resume certain systems or timers
}

func (mms *MainMenuScene) AddEntity(entity *ecs.Entity) {
	mms.entities = append(mms.entities, entity)
}

func (mms *MainMenuScene) RemoveAllEntities() {
	for _, entity := range mms.entities {
		mms.ecsManager.DestroyEntity(entity.ID)
	}
	mms.entities = []*ecs.Entity{}
}
