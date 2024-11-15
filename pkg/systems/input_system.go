package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/editor"
	physicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"
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
	//is.handlePlayerMovementInput()
	is.handlePlayerMovmentInput()

	// Handle editor input
	is.handleEditorInput()

	// Handle spawner input
	is.handleSpawnerInput()

	// Handle rigid body spawner
	is.handleRigidBodySpawner()
}

func (is *InputSystem) handlePlayerMovmentInput() {

	playerComps, playerCompsExists := is.componentsManager.Components[ecs.PlayerComponent]

	if !playerCompsExists {
		utils.ErrorLogger.Println("InputSystem: No player components found")
		return
	}

	for entity := range playerComps {
		rb, rbExists := is.EcsManager.GetComponent(entity, ecs.RigidBodyComponent).(*physicscomponents.RigidBody)

		if !rbExists {
			continue
		}

		// Define the movement force
		movementForce := float32(2000.0)

		// Initialize the movement vector
		force := raylib.NewVector2(0, 0)

		if raylib.IsKeyDown(raylib.KeyW) {
			force.Y = -movementForce
		}

		if raylib.IsKeyDown(raylib.KeyS) {
			force.Y = movementForce
		}

		if raylib.IsKeyDown(raylib.KeyA) {
			force.X = -movementForce
		}

		if raylib.IsKeyDown(raylib.KeyD) {
			force.X = movementForce
		}

		// Apply the movement force to the RigidBody's force
		rb.Force = raylib.Vector2Add(rb.Force, force)

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
			is.EcsManager.AddComponent(entity, ecs.Transform2DComponent, &components.Transform2D{Position: raylib.NewVector2(float32(raylib.GetRandomValue(0, int32(raylib.GetScreenWidth())-1)), float32(raylib.GetRandomValue(0, int32(raylib.GetScreenHeight())-1)))})
			is.EcsManager.AddComponent(entity, ecs.VelocityComponent, &components.Velocity{Vector: raylib.NewVector2(float32(raylib.GetRandomValue(-10, 10)), float32(raylib.GetRandomValue(-10, 10)))})
			is.EcsManager.AddComponent(entity, ecs.SpeedComponent, &components.Speed{Value: float32(raylib.GetRandomValue(50, 200))})
			is.EcsManager.AddComponent(entity, ecs.ColorComponent, &components.Color{Color: color})
		}
	}
}

func (is *InputSystem) handleRigidBodySpawner() {

	if raylib.IsKeyPressed(raylib.KeyR) {

		// Create a rigid body entity
		rigidBodyEntity := is.EcsManager.CreateEntity()

		// Add a rigid body component to the rigid body entity
		is.EcsManager.AddComponent(rigidBodyEntity, ecs.RigidBodyComponent, &physicscomponents.RigidBody{
			Mass:         1,
			Velocity:     raylib.NewVector2(0, 0),
			Acceleration: raylib.NewVector2(0, 0),
			Force:        raylib.NewVector2(0, 0),
			IsKinematic:  false,
			IsStatic:     false,
		})

	}
}
