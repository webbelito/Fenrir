package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/ecs"
)

type InputSystem struct{}

func (is *InputSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	// Example for single player entity
	playerEntity := ecs.Entity(1)

	// Get the velocity component for the player entity
	velocityComp, velExists := cm.Components[ecs.VelocityComponent][playerEntity].(*ecs.Velocity)

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
