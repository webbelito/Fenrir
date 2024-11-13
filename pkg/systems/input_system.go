package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/editor"
	"github.com/webbelito/Fenrir/pkg/utils"
)

type InputSystem struct {
	EcsManager        *ecs.ECSManager
	Editor            *editor.Editor
	entitiesManager   *ecs.EntitiesManager
	componentsManager *ecs.ComponentsManager
}

func (is *InputSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	if em == nil || cm == nil || is.Editor == nil {
		utils.ErrorLogger.Println("InputSystem: EntitiesManager, ComponentsManager or Editor is nil")
		return
	}

	// Assign the entities and components manager to the system
	is.entitiesManager = em
	is.componentsManager = cm

	// Handle player movement input
	is.handlePlayerMovementInput()

	// Handle editor input
	is.handleEditorInput()

	// Handle spawner input
	is.handleSpawnerInput()
}

func (is *InputSystem) handlePlayerMovementInput() {

	// Get all the player components
	playerComps, playerCompsExists := is.componentsManager.Components[ecs.PlayerComponent]

	if !playerCompsExists {
		utils.ErrorLogger.Println("InputSystem: No player components found")
		return
	}

	playerEntity := &ecs.Entity{}

	// Get the first player component from playerComps
	for player := range playerComps {
		playerEntity = player
		break
	}

	// Get the velocity component for the player entity
	velocityComp, velExists := is.componentsManager.Components[ecs.VelocityComponent][playerEntity].(*components.Velocity)

	if !velExists {
		return
	}

	// Reset the velocity to zero; it will be set based on input
	velocityComp.Vector = raylib.NewVector2(0, 0)

	// Apply input-based acceleration
	if raylib.IsKeyDown(raylib.KeyW) {
		velocityComp.Vector.Y = -1
	}
	if raylib.IsKeyDown(raylib.KeyS) {
		velocityComp.Vector.Y = 1
	}
	if raylib.IsKeyDown(raylib.KeyA) {
		velocityComp.Vector.X = -1
	}
	if raylib.IsKeyDown(raylib.KeyD) {
		velocityComp.Vector.X = 1
	}

}

func (is *InputSystem) handleEditorInput() {
	if raylib.IsKeyPressed(raylib.KeyF1) {
		is.Editor.ToggleVisibility()
	}
}

// TODO: Remove this function as it's more of a test function
func (is *InputSystem) handleSpawnerInput() {

	// Colors to choose from
	colors := []raylib.Color{
		raylib.Blue,
		raylib.Green,
		raylib.Purple,
		raylib.Orange,
		raylib.Pink,
		raylib.Yellow,
		raylib.SkyBlue,
		raylib.Lime,
		raylib.Gold,
		raylib.Violet,
		raylib.Brown,
		raylib.LightGray,
		raylib.DarkGray,
	}

	if raylib.IsKeyPressed(raylib.KeySpace) {

		// Create 500 entities with random positions, velocities, speeds and colors
		for i := 0; i < 500; i++ {

			// Select a random color from the colors slice
			color := colors[raylib.GetRandomValue(0, int32(len(colors)-1))]

			// Create an entity with a random position, velocity, speed and color
			entity := is.EcsManager.CreateEntity()
			is.EcsManager.AddComponent(entity, ecs.PositionComponent, &components.Position{Vector: raylib.NewVector2(float32(raylib.GetRandomValue(0, int32(raylib.GetScreenWidth())-1)), float32(raylib.GetRandomValue(0, int32(raylib.GetScreenHeight())-1)))})
			is.EcsManager.AddComponent(entity, ecs.VelocityComponent, &components.Velocity{Vector: raylib.NewVector2(float32(raylib.GetRandomValue(-10, 10)), float32(raylib.GetRandomValue(-10, 10)))})
			is.EcsManager.AddComponent(entity, ecs.SpeedComponent, &components.Speed{Value: float32(raylib.GetRandomValue(50, 200))})
			is.EcsManager.AddComponent(entity, ecs.ColorComponent, &components.Color{Color: color})
		}
	}
}
