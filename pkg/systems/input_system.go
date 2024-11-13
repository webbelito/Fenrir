package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"
)

type InputSystem struct {
	entitiesManager   *ecs.EntitiesManager
	componentsManager *ecs.ComponentsManager
}

func (is *InputSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	if em == nil || cm == nil {
		utils.ErrorLogger.Println("InputSystem: EntitiesManager or ComponentsManager is nil")
		return
	}

	// Assign the entities and components manager to the system
	is.entitiesManager = em
	is.componentsManager = cm

	// Handle player movement input
	is.handlePlayerMovementInput()

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
