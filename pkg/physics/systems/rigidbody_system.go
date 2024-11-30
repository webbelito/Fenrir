package physicssystems

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	physicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// RigidBodySystem is a system that handles rigid body physics
type RigidBodySystem struct {
	gravity  raylib.Vector2
	manager  *ecs.Manager
	priority int
}

// NewRigidBodySystem creates a new RigidBodySystem
func NewRigidBodySystem(m *ecs.Manager, g raylib.Vector2, p int) *RigidBodySystem {
	return &RigidBodySystem{
		gravity:  g,
		manager:  m,
		priority: p,
	}
}

func (rbs *RigidBodySystem) Update(dt float64) {

	// Get all entities with RigidBodyComponent
	entities := rbs.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.RigidBodyComponent})

	// Get the screen width and height
	screenWidth := float32(raylib.GetScreenWidth())
	screenHeight := float32(raylib.GetScreenHeight())

	// Iterate over all entities with RigidBodyComponent
	for _, entity := range entities {

		// Get the RigidBodyComponent
		rigidBodyComp, exist := rbs.manager.GetComponent(entity, ecs.RigidBodyComponent)

		// Check if the RigidBodyComponent exists
		if !exist {
			utils.ErrorLogger.Println("RigidBodyComponent does not exist")
			continue
		}

		// Cast the component to a RigidBody
		rb, ok := rigidBodyComp.(*physicscomponents.RigidBody)

		// Check if the cast was successful
		if !ok {
			utils.ErrorLogger.Println("Failed to cast RigidBodyComponent to RigidBody")
			continue
		}

		if rb.IsStatic {
			continue
		}

		if !rb.IsKinematic {
			// Apply gravity to the entity (F = m * g)
			rb.Force = raylib.Vector2Add(rb.Force, raylib.Vector2Scale(rbs.gravity, rb.Mass))
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

		// Calculate InvMass
		if rb.Mass != 0 {
			rb.InvMass = 1 / rb.Mass
		} else {
			rb.InvMass = 0
		}

		// Update velocity based on acceleration (v += a * dt)
		rb.Velocity = raylib.Vector2Add(rb.Velocity, raylib.Vector2Scale(rb.Acceleration, float32(dt)))

		// Get the transform component for the entity
		transformComp, exist := rbs.manager.GetComponent(entity, ecs.Transform2DComponent)

		// Check if the Transform2DComponent exists
		if !exist {
			utils.ErrorLogger.Println("Transform2DComponent does not exist")
			continue
		}

		// Cast the component to a Transform2D
		transform, ok := transformComp.(*components.Transform2D)

		// Check if the cast was successful
		if !ok {
			utils.ErrorLogger.Println("Failed to cast Transform2DComponent to Transform2D")
			continue
		}

		// Clamp the position to the screen bounds

		// If the position is less than the size, set it to the size
		if transform.Position.X < transform.Scale.X {
			transform.Position.X = transform.Scale.X
			rb.Velocity.X = 0
			// If the position is greater than the screen width minus the size, set it to the screen width minus the size
		} else if transform.Position.X > screenWidth-transform.Scale.X {
			transform.Position.X = screenWidth - transform.Scale.X
			rb.Velocity.X = 0
		}

		// If the position is less than the size, set it to the size
		if transform.Position.Y < transform.Scale.Y {
			transform.Position.Y = transform.Scale.Y
			rb.Velocity.Y = 0
			// If the position is greater than the screen height minus the size, set it to the screen height minus the size
		} else if transform.Position.Y > screenHeight-transform.Scale.Y {
			transform.Position.Y = screenHeight - transform.Scale.Y
			rb.Velocity.Y = 0
		}

		// Update position based on velocity (p += v * dt)
		transform.Position = raylib.Vector2Add(transform.Position, raylib.Vector2Scale(rb.Velocity, float32(dt)))

		// Reset force for the next frame
		rb.Force = raylib.NewVector2(0, 0)
	}
}

// GetPriority returns the priority of the system
func (rbs *RigidBodySystem) GetPriority() int {
	return rbs.priority
}
