package scenes

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/events"
	"github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type MainMenuScene struct {
	sceneManager *SceneManager
	ecsManager   *ecs.ECSManager
	sceneData    *SceneData

	entities        []*ecs.Entity
	logicSystems    []systeminterfaces.UpdatableSystemInterface
	renderSystems   []systeminterfaces.RenderableSystemInterface
	uiRenderSystems []systeminterfaces.UIRenderableSystemInterface
}

func NewMainMenuScene(sm *SceneManager, em *ecs.ECSManager, sd *SceneData) *MainMenuScene {
	mms := &MainMenuScene{
		sceneManager:    sm,
		ecsManager:      em,
		sceneData:       sd,
		entities:        []*ecs.Entity{},
		logicSystems:    []systeminterfaces.UpdatableSystemInterface{},
		renderSystems:   []systeminterfaces.RenderableSystemInterface{},
		uiRenderSystems: []systeminterfaces.UIRenderableSystemInterface{},
	}

	return mms

}

// TODO: Immplement Init function
func (mms *MainMenuScene) Initialize() {

	utils.InfoLogger.Println("Initializing Main Menu Scene...")

	mms.initializeUIEntities()

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
		mms.ecsManager.GetEventsManager().Dispatch("exit_game", events.ExitGameEvent{ShouldExitGame: true})
	}
}

func (mms *MainMenuScene) Render() {

	// Currently only rendering the UI with the UI Render System

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

func (mms *MainMenuScene) initializeUIEntities() {
	for _, entityData := range mms.sceneData.Entities {
		entity := mms.ecsManager.CreateEntity()

		for componentType, componentData := range entityData.Components {
			switch componentType {
			case "UIPanel":
				panel, panelOk := componentData.(map[string]interface{})
				if !panelOk {
					utils.ErrorLogger.Println("Failed to assert component data as map[string]interface{}")
					continue
				}

				mms.addUIPanel(entity, panel)

			case "UIButton":
				button, buttonOk := componentData.(map[string]interface{})
				if !buttonOk {
					utils.ErrorLogger.Println("Failed to assert component data as map[string]interface{}")
					continue
				}

				mms.addUIButton(entity, button)

			case "UILabel":
				label, labelOk := componentData.(map[string]interface{})
				if !labelOk {
					utils.ErrorLogger.Println("Failed to assert component data as []interface{}")
					continue
				}

				mms.addUILabel(entity, label)

			default:
				utils.ErrorLogger.Println("Unknown component type: ", componentType)
			}
		}
	}
}

func (mms *MainMenuScene) addUIPanel(e *ecs.Entity, data map[string]interface{}) {
	title, ok := data["Title"].(string)
	if !ok {
		utils.ErrorLogger.Println("Invalid Title for UIPanel")
		return
	}

	boundsMap, ok := data["Bounds"].(map[string]interface{})
	if !ok {
		utils.ErrorLogger.Println("Invalid Bounds for UIPanel")
		return
	}

	bounds := raylib.Rectangle{
		X:      float32(boundsMap["X"].(float64)),
		Y:      float32(boundsMap["Y"].(float64)),
		Width:  float32(boundsMap["Width"].(float64)),
		Height: float32(boundsMap["Height"].(float64)),
	}

	visible, ok := data["IsVisible"].(bool)
	if !ok {
		utils.ErrorLogger.Println("Invalid IsVisible for UIPanel")
		return
	}

	panel := &components.UIPanel{
		Title:     title,
		Bounds:    bounds,
		IsVisible: visible,
	}

	mms.ecsManager.GetUIComponentsManager().AddComponent(e.ID, ecs.UIPanelComponent, panel)
	mms.AddEntity(e)
}

func (mms *MainMenuScene) addUIButton(e *ecs.Entity, data map[string]interface{}) {
	text, ok := data["Text"].(string)
	if !ok {
		utils.ErrorLogger.Println("Invalid Text for UIButton")
		return
	}

	boundsMap, ok := data["Bounds"].(map[string]interface{})
	if !ok {
		utils.ErrorLogger.Println("Invalid Bounds for UIButton")
		return
	}

	bounds := raylib.Rectangle{
		X:      float32(boundsMap["X"].(float64)),
		Y:      float32(boundsMap["Y"].(float64)),
		Width:  float32(boundsMap["Width"].(float64)),
		Height: float32(boundsMap["Height"].(float64)),
	}

	visible, ok := data["IsVisible"].(bool)
	if !ok {
		utils.ErrorLogger.Println("Invalid IsVisible for UIButton")
		return
	}

	button := &components.UIButton{
		Text:      text,
		Bounds:    bounds,
		IsVisible: visible,
	}

	mms.ecsManager.GetUIComponentsManager().AddComponent(e.ID, ecs.UIButtonComponent, button)
	mms.AddEntity(e)
}

func (mms *MainMenuScene) addUILabel(e *ecs.Entity, data map[string]interface{}) {

	label, ok := data["Label"].(string)
	if !ok {
		utils.ErrorLogger.Println("Invalid Label for UILabel")
		return
	}

	boundsMap, ok := data["Bounds"].(map[string]interface{})

	if !ok {
		utils.ErrorLogger.Println("Invalid Bounds for UILabel")
		return
	}

	bounds := raylib.Rectangle{
		X:      float32(boundsMap["X"].(float64)),
		Y:      float32(boundsMap["Y"].(float64)),
		Width:  float32(boundsMap["Width"].(float64)),
		Height: float32(boundsMap["Height"].(float64)),
	}

	visible, ok := data["IsVisible"].(bool)
	if !ok {
		utils.ErrorLogger.Println("Invalid IsVisible for UILabel")
		return
	}

	labelComp := &components.UILabel{
		Label:     label,
		Bounds:    bounds,
		IsVisible: visible,
	}

	mms.ecsManager.GetUIComponentsManager().AddComponent(e.ID, ecs.UILabelComponent, labelComp)
	mms.AddEntity(e)
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
