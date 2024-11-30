package systems

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// MovementSystem is a system that handles movement
type MovementSystem struct {
	manager  *ecs.Manager
	priority int
}

// NewMovementSystem creates a new MovementSystem
func NewMovementSystem(m *ecs.Manager, p int) *MovementSystem {
	return &MovementSystem{
		manager:  m,
		priority: p,
	}
}

// Update updates the MovementSystem
func (ms *MovementSystem) Update(dt float64) {

	// Check if the ECS Manager is nil
	if ms.manager == nil {
		utils.ErrorLogger.Println("MovementSystem: ECS Manager is nil")
		return
	}

	// Move entities
	ms.MoveEntities(dt)

}

// MoveEntities moves entities
func (ms *MovementSystem) MoveEntities(dt float64) {

	entityIDs := ms.manager.GetEntitiesWithComponents([]ecs.ComponentType{
		ecs.Transform2DComponent,
		ecs.VelocityComponent,
		ecs.SpeedComponent,
	})

	// If there are no entities with the required components, log error and return
	if len(entityIDs) == 0 {
		utils.ErrorLogger.Println("MovementSystem: No entities with required components found")
		return
	}

	// Iterate over all entities with the required components
	for _, eID := range entityIDs {

		// Get the transform component
		transformComp, exist := ms.manager.GetComponent(eID, ecs.Transform2DComponent)

		if !exist {
			utils.ErrorLogger.Println("MovementSystem: Transform2D component does not exist")
			return
		}

		// Cast the transform component to the correct type
		transform, ok := transformComp.(*components.Transform2D)

		if !ok {
			utils.ErrorLogger.Println("MovementSystem: Transform2D component could not be casted")
			return
		}

		// Get the velocity component
		velocityComp, exist := ms.manager.GetComponent(eID, ecs.VelocityComponent)

		if !exist {
			utils.ErrorLogger.Println("MovementSystem: Velocity component does not exist")
			return
		}

		// Cast the velocity component to the correct type
		velocity, ok := velocityComp.(*components.Velocity)

		if !ok {
			utils.ErrorLogger.Println("MovementSystem: Velocity component could not be casted")
			return
		}

		// Get the speed component
		speedComp, exist := ms.manager.GetComponent(eID, ecs.SpeedComponent)

		if !exist {
			utils.ErrorLogger.Println("MovementSystem: Speed component does not exist")
			return
		}

		// Cast the speed component to the correct type
		speed, ok := speedComp.(*components.Speed)

		if !ok {
			utils.ErrorLogger.Println("MovementSystem: Speed component could not be casted")
			return
		}

		// Normalize the velocity vector to ensure consistent movement speed
		normalizedVelocity := raylib.Vector2Normalize(velocity.Vector)

		// Calculate the new position based on the velocity and speed
		deltaVelocity := raylib.Vector2Scale(normalizedVelocity, speed.Value*float32(dt))

		// Update the position of the entity
		transform.Position = raylib.Vector2Add(transform.Position, deltaVelocity)
	}
}

func (ms *MovementSystem) GetPriority() int {
	return ms.priority
}
