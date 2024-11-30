package scenes

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/events"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// MainMenuScene is a struct that represents the Main Menu Scene
type MainMenuScene struct {
	sceneManager *SceneManager
	manager      *ecs.Manager
	sceneData    *SceneData

	sceneEntities []*ecs.Entity
}

// NewMainMenuScene creates a new MainMenuScene
func NewMainMenuScene(sm *SceneManager, m *ecs.Manager, sd *SceneData) *MainMenuScene {
	return &MainMenuScene{
		sceneManager:  sm,
		manager:       m,
		sceneData:     sd,
		sceneEntities: []*ecs.Entity{},
	}
}

// Initialize initializes the Main Menu Scene
func (mms *MainMenuScene) Initialize() {

	mms.initializeUIEntities()
}

func (mms *MainMenuScene) Update(dt float64) {
	// Check if the ECS Manager is nil
	if mms.manager == nil {
		utils.ErrorLogger.Println("MainMenuScene: ECS Manager is nil")
		return
	}

	// Handle Main Menu Input
	if raylib.IsKeyPressed(raylib.KeyEnter) {

		// Change to Game Scene
		err := mms.sceneManager.ChangeScene("assets/scenes/game_scene.json")
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}

	// TODO: Remove temporary input handling
	if raylib.IsKeyPressed(raylib.KeyEscape) {
		mms.manager.DispatchEvent("exit_game", events.ExitGameEvent{ShouldExitGame: true})
	}
}

func (mms *MainMenuScene) Render() {
	// Currently only rendering the UI with the UI Render System
}

// Cleanup cleans up the Main Menu Scene
func (mms *MainMenuScene) Cleanup() {
	// Remove all entities created by this scene
	mms.RemoveAllEntities()
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

// initializeUIEntities initializes the UI entities for the Main Menu Scene
func (mms *MainMenuScene) initializeUIEntities() {

	// Loop through the entities in the scene data
	for _, entityData := range mms.sceneData.Entities {

		// Create a new entity
		entity := mms.manager.CreateEntity()

		// Loop through the components in the entity data
		for componentType, componentData := range entityData.Components {

			switch componentType {
			case "UIPanel":
				panel, ok := componentData.(map[string]interface{})
				if !ok {
					utils.ErrorLogger.Println("Failed to assert component data as map[string]interface{}")
					continue
				}

				// Add UIPanel component to the entity
				mms.addUIPanel(entity, panel)

			case "UIButton":
				button, ok := componentData.(map[string]interface{})
				if !ok {
					utils.ErrorLogger.Println("Failed to assert component data as map[string]interface{}")
					continue
				}

				// Add UIButton component to the entity
				mms.addUIButton(entity, button)

			case "UILabel":
				label, ok := componentData.(map[string]interface{})
				if !ok {
					utils.ErrorLogger.Println("Failed to assert component data as []interface{}")
					continue
				}

				// Add UILabel component to the entity
				mms.addUILabel(entity, label)

			default:
				utils.ErrorLogger.Println("Unknown component type: ", componentType)
			}
		}
	}
}

// addUIPanel adds a UIPanel component to an entity
func (mms *MainMenuScene) addUIPanel(e *ecs.Entity, data map[string]interface{}) {

	// Get the title of the panel
	title, ok := data["Title"].(string)

	// Check if the title is valid
	if !ok {
		utils.ErrorLogger.Println("Invalid Title for UIPanel")
		return
	}

	// Get the bounds of the panel
	boundsMap, ok := data["Bounds"].(map[string]interface{})

	// Check if the bounds are valid
	if !ok {
		utils.ErrorLogger.Println("Invalid Bounds for UIPanel")
		return
	}

	// Create a new rectangle for the bounds
	bounds := raylib.Rectangle{
		X:      float32(boundsMap["X"].(float64)),
		Y:      float32(boundsMap["Y"].(float64)),
		Width:  float32(boundsMap["Width"].(float64)),
		Height: float32(boundsMap["Height"].(float64)),
	}

	// Get the visibility of the panel
	visible, ok := data["IsVisible"].(bool)

	// Check if the visibility is valid
	if !ok {
		utils.ErrorLogger.Println("Invalid IsVisible for UIPanel")
		return
	}

	// Create a new UIPanel component
	panel := &components.UIPanel{
		Title:     title,
		Bounds:    bounds,
		IsVisible: visible,
	}

	// Add the UIPanel component to the entity
	mms.manager.AddComponent(e.ID, ecs.UIPanelComponent, panel)

	// Add the entity to the scene
	mms.AddEntity(e)
}

// addUIButton adds a UIButton component to an entity
func (mms *MainMenuScene) addUIButton(e *ecs.Entity, data map[string]interface{}) {

	// Get the text of the button
	text, ok := data["Text"].(string)

	// Check if the text is valid
	if !ok {
		utils.ErrorLogger.Println("Invalid Text for UIButton")
		return
	}

	// Get the bounds of the button
	boundsMap, ok := data["Bounds"].(map[string]interface{})

	// Check if the bounds are valid
	if !ok {
		utils.ErrorLogger.Println("Invalid Bounds for UIButton")
		return
	}

	// Create a new rectangle for the bounds
	bounds := raylib.Rectangle{
		X:      float32(boundsMap["X"].(float64)),
		Y:      float32(boundsMap["Y"].(float64)),
		Width:  float32(boundsMap["Width"].(float64)),
		Height: float32(boundsMap["Height"].(float64)),
	}

	// Get the visibility of the button
	visible, ok := data["IsVisible"].(bool)

	// Check if the visibility is valid
	if !ok {
		utils.ErrorLogger.Println("Invalid IsVisible for UIButton")
		return
	}

	// Create a new UIButton component
	button := &components.UIButton{
		Text:      text,
		Bounds:    bounds,
		IsVisible: visible,
	}

	// Add the UIButton component to the entity
	mms.manager.AddComponent(e.ID, ecs.UIButtonComponent, button)

	// Add the entity to the scene
	mms.AddEntity(e)
}

// addUILabel adds a UILabel component to an entity
func (mms *MainMenuScene) addUILabel(e *ecs.Entity, data map[string]interface{}) {

	// Get the label text
	label, ok := data["Label"].(string)

	// Check if the label is valid
	if !ok {
		utils.ErrorLogger.Println("Invalid Label for UILabel")
		return
	}

	// Get the bounds of the label
	boundsMap, ok := data["Bounds"].(map[string]interface{})

	// Check if the bounds are valid
	if !ok {
		utils.ErrorLogger.Println("Invalid Bounds for UILabel")
		return
	}

	// Create a new rectangle for the bounds
	bounds := raylib.Rectangle{
		X:      float32(boundsMap["X"].(float64)),
		Y:      float32(boundsMap["Y"].(float64)),
		Width:  float32(boundsMap["Width"].(float64)),
		Height: float32(boundsMap["Height"].(float64)),
	}

	// Get the visibility of the label
	visible, ok := data["IsVisible"].(bool)

	// Check if the visibility is valid
	if !ok {
		utils.ErrorLogger.Println("Invalid IsVisible for UILabel")
		return
	}

	// Create a new UILabel component
	labelComp := &components.UILabel{
		Label:     label,
		Bounds:    bounds,
		IsVisible: visible,
	}

	// Add the UILabel component to the entity
	mms.manager.AddComponent(e.ID, ecs.UILabelComponent, labelComp)

	// Add the entity to the scene
	mms.AddEntity(e)
}

// AddEntity adds an entity to the scene
func (mms *MainMenuScene) AddEntity(entity *ecs.Entity) {
	mms.sceneEntities = append(mms.sceneEntities, entity)
}

// RemoveAllEntities removes all entities created by the scene
func (mms *MainMenuScene) RemoveAllEntities() {

	// Iterate over all entities created by the scene
	for _, entity := range mms.sceneEntities {
		mms.manager.DestroyEntity(entity.ID)
	}

	// Clear the scene entities
	mms.sceneEntities = []*ecs.Entity{}
}
