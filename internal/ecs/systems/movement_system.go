// internal/ecs/systems/movement_system.go

package systems

import (
	"fmt"

	"github.com/webbelito/Fenrir/internal/ecs"
	"github.com/webbelito/Fenrir/internal/ecs/components"
)

type MovementSystem struct{}

func (ms *MovementSystem) Update(deltaTime float32, manager *ecs.Manager) {

	// Get all entities with a position and velocity component
	entities := manager.GetEntitiesWithComponents(&components.Position{}, &components.Velocity{})

	for _, entity := range entities {
		position := manager.GetComponent(entity, &components.Position{}).(*components.Position)
		velocity := manager.GetComponent(entity, &components.Velocity{}).(*components.Velocity)

		// Move the entity
		position.X += velocity.X * deltaTime
		position.Y += velocity.Y * deltaTime

		fmt.Printf("Entity %d moved to (%f, %f)\n", entity, position.X, position.Y)
	}

}
