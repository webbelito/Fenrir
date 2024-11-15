package physicssystems

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	physicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type RigidBodySystem struct {
	Gravity raylib.Vector2
}

func NewRigidBodySystem(gravity raylib.Vector2) *RigidBodySystem {
	return &RigidBodySystem{
		Gravity: gravity,
	}
}

func (rbs *RigidBodySystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	rigidBodyComps, rigidBodyCompsExists := cm.Components[ecs.RigidBodyComponent]

	if !rigidBodyCompsExists {
		return
	}

	screenWidth := float32(raylib.GetScreenWidth())
	screenHeight := float32(raylib.GetScreenHeight())

	for entity, rigidBodyComp := range rigidBodyComps {
		rb, rbExists := rigidBodyComp.(*physicscomponents.RigidBody)

		if !rbExists {
			continue
		}

		if rb.IsStatic {
			continue
		}

		if !rb.IsKinematic {
			// Apply gravity to the entity (F = m * g)
			rb.Force = raylib.Vector2Add(rb.Force, raylib.Vector2Scale(rbs.Gravity, rb.Mass))
		}

		// Apply drag: F Drag = -Drag * v
		dragForce := raylib.Vector2Scale(rb.Velocity, -rb.Drag)
		rb.Force = raylib.Vector2Add(rb.Force, dragForce)

		// Update acceleration based on force (a = F / m)
		if rb.Mass != 0 {
			rb.Acceleration = raylib.Vector2Scale(rb.Force, 1/rb.Mass)
		} else {
			rb.Acceleration = raylib.NewVector2(0, 0)
		}

		// Update velocity based on acceleration (v += a * dt)
		rb.Velocity = raylib.Vector2Add(rb.Velocity, raylib.Vector2Scale(rb.Acceleration, float32(dt)))

		// Get the position component for the entity
		positionComp, posCompExists := cm.Components[ecs.PositionComponent][entity].(*components.Position)

		// Handle Position related Updates
		if posCompExists {

			// Boundry clamping
			// Assume entites have a Size component
			sizeComp, sizeCompExists := cm.Components[ecs.SizeComponent][entity].(*components.Size)

			if sizeCompExists {

				// Clamp the position to the screen bounds

				// If the position is less than the size, set it to the size
				if positionComp.Vector.X < sizeComp.Size.X {
					positionComp.Vector.X = sizeComp.Size.X
					rb.Velocity.X = 0
					// If the position is greater than the screen width minus the size, set it to the screen width minus the size
				} else if positionComp.Vector.X > screenWidth-sizeComp.Size.X {
					positionComp.Vector.X = screenWidth - sizeComp.Size.X
					rb.Velocity.X = 0
				}

				// If the position is less than the size, set it to the size
				if positionComp.Vector.Y < sizeComp.Size.Y {
					positionComp.Vector.Y = sizeComp.Size.Y
					rb.Velocity.Y = 0
					// If the position is greater than the screen height minus the size, set it to the screen height minus the size
				} else if positionComp.Vector.Y > screenHeight-sizeComp.Size.Y {
					positionComp.Vector.Y = screenHeight - sizeComp.Size.Y
					rb.Velocity.Y = 0
				}
			}

			// Update position based on velocity (p += v * dt)
			positionComp.Vector = raylib.Vector2Add(positionComp.Vector, raylib.Vector2Scale(rb.Velocity, float32(dt)))
		}

		// Reset force for the next frame
		rb.Force = raylib.NewVector2(0, 0)

	}
}
