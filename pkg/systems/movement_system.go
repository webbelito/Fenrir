package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"
)

type MovementSystem struct {
	entitiesManager   *ecs.EntitiesManager
	componentsManager *ecs.ComponentsManager
}

func NewMovementSystem() *MovementSystem {
	return &MovementSystem{}
}

func (ms *MovementSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	if em == nil || cm == nil {
		utils.ErrorLogger.Println("MovementSystem: EntitiesManager or ComponentsManager is nil")
		return
	}

	// Assign the entities and components manager to the system
	ms.entitiesManager = em
	ms.componentsManager = cm

	// Move entities
	ms.MoveEntities(dt)

}

func (ms *MovementSystem) MoveEntities(dt float64) {

	// Get entities with a position, velocity, speed, component
	positionComps, posExist := ms.componentsManager.Components[ecs.PositionComponent]
	veloComps, velExist := ms.componentsManager.Components[ecs.VelocityComponent]
	SpeedComps, speedExist := ms.componentsManager.Components[ecs.SpeedComponent]
	PlayerComps := ms.componentsManager.Components[ecs.PlayerComponent]

	if !posExist || !velExist || !speedExist {
		return
	}

	// Update the position of all entities with a position, velocity and speed component
	for entity, vel := range veloComps {
		position, posExists := positionComps[entity].(*components.Position)
		velocity, velExists := vel.(*components.Velocity)
		speed, speedExists := SpeedComps[entity].(*components.Speed)
		_, playerExists := PlayerComps[entity].(*components.Player)

		if !posExists || !velExists || !speedExists {
			continue
		}

		// Normalize the velocity vector to ensure consistent movement speed
		normalizedVelocity := raylib.Vector2Normalize(velocity.Vector)

		// Calculate the new position based on the velocity and speed
		deltaVelocity := raylib.Vector2Scale(normalizedVelocity, speed.Value*float32(dt))
		position.Vector = raylib.Vector2Add(position.Vector, deltaVelocity)

		// Clamp the player position to the screen bounds
		if playerExists {

			// Define the screen bounds
			screenWidth := float32(raylib.GetScreenWidth())
			screenHeight := float32(raylib.GetScreenHeight())

			// Clamp the position to the screen bounds
			position.Vector.X = raylib.Clamp(position.Vector.X, 0, screenWidth-5)
			position.Vector.Y = raylib.Clamp(position.Vector.Y, 0, screenHeight-5)

		}
	}
}
