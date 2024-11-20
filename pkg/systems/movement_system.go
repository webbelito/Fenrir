package systems

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type MovementSystem struct {
	ecsManager        *ecs.ECSManager
	entitiesManager   *ecs.EntitiesManager
	componentsManager *ecs.ComponentsManager
	priority          int
}

func NewMovementSystem(ecsM *ecs.ECSManager, p int) *MovementSystem {
	return &MovementSystem{
		ecsManager:        ecsM,
		entitiesManager:   ecsM.GetEntitiesManager(),
		componentsManager: ecsM.GetComponentsManager(),
		priority:          p,
	}
}

func (ms *MovementSystem) Update(dt float64) {

	if ms.ecsManager == nil || ms.entitiesManager == nil || ms.componentsManager == nil {
		utils.ErrorLogger.Println("MovementSystem: ECSManager or EntitiesManager or ComponentsManager is nil")
		return
	}

	// Move entities
	ms.MoveEntities(dt)

}

func (ms *MovementSystem) MoveEntities(dt float64) {

	// Get entities with a position, velocity, speed, component
	transformComps, tCompsExist := ms.componentsManager.Components[ecs.Transform2DComponent]
	veloComps, vCompsExists := ms.componentsManager.Components[ecs.VelocityComponent]
	SpeedComps, sCompsExists := ms.componentsManager.Components[ecs.SpeedComponent]
	PlayerComps := ms.componentsManager.Components[ecs.PlayerComponent]

	if !tCompsExist || !vCompsExists || !sCompsExists {
		return
	}

	// Update the position of all entities with a position, velocity and speed component
	for entity, vel := range veloComps {
		transform, tExists := transformComps[entity].(*components.Transform2D)
		velocity, vExists := vel.(*components.Velocity)
		speed, sExists := SpeedComps[entity].(*components.Speed)
		_, pExists := PlayerComps[entity].(*components.Player)

		if !tExists || !vExists || !sExists {
			continue
		}

		// Normalize the velocity vector to ensure consistent movement speed
		normalizedVelocity := raylib.Vector2Normalize(velocity.Vector)

		// Calculate the new position based on the velocity and speed
		deltaVelocity := raylib.Vector2Scale(normalizedVelocity, speed.Value*float32(dt))
		transform.Position = raylib.Vector2Add(transform.Position, deltaVelocity)

		// Clamp the player position to the screen bounds
		if pExists {

			// Define the screen bounds
			screenWidth := float32(raylib.GetScreenWidth())
			screenHeight := float32(raylib.GetScreenHeight())

			// Clamp the position to the screen bounds
			transform.Position.X = raylib.Clamp(transform.Position.X, 0, screenWidth-5)
			transform.Position.Y = raylib.Clamp(transform.Position.Y, 0, screenHeight-5)

		}
	}
}

func (ms *MovementSystem) GetPriority() int {
	return ms.priority
}
