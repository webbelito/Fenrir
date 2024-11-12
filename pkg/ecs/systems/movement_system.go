package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/ecs"
)

type MovementSystem struct{}

func (ms *MovementSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	// Get entities with a position, velocity, speed, component
	positionComps, posExist := cm.Components[ecs.PositionComponent]
	veloComps, velExist := cm.Components[ecs.VelocityComponent]
	SpeedComps, speedExist := cm.Components[ecs.SpeedComponent]

	if !posExist || !velExist || !speedExist {
		return
	}

	// Update the position of all entities with a position, velocity and speed component
	for entity, vel := range veloComps {
		position, posExists := positionComps[entity].(*ecs.Position)
		velocity, velExists := vel.(*ecs.Velocity)
		speed, speedExists := SpeedComps[entity].(*ecs.Speed)

		if !posExists || !velExists || !speedExists {
			continue
		}

		// Normalize the velocity vector to ensure consistent movement speed
		normalizedVelocity := raylib.Vector2Normalize(velocity.Vector)

		// Calculate the new position based on the velocity and speed
		deltaVelocity := raylib.Vector2Scale(normalizedVelocity, speed.Value*float32(dt))
		position.Vector = raylib.Vector2Add(position.Vector, deltaVelocity)

		// Define the screen bounds
		screenWidth := float32(raylib.GetScreenWidth())
		screenHeight := float32(raylib.GetScreenHeight())

		// Clamp the position to the screen bounds
		position.Vector.X = raylib.Clamp(position.Vector.X, 0, screenWidth)
		position.Vector.Y = raylib.Clamp(position.Vector.Y, 0, screenHeight)

	}
}
