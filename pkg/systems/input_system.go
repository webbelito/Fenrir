package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/editor"
	physicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"
	"github.com/webbelito/Fenrir/pkg/utils"
)

// InputSystem is a system that handles input
type InputSystem struct {
	manager  *ecs.Manager
	editor   *editor.EditorManager
	priority int
}

// NewInputSystem creates a new InputSystem
func NewInputSystem(m *ecs.Manager, e *editor.EditorManager, p int) *InputSystem {
	return &InputSystem{
		manager:  m,
		editor:   e,
		priority: p,
	}
}

func (is *InputSystem) Update(dt float64) {

	if is.manager == nil || is.editor == nil {
		utils.ErrorLogger.Println("InputSystem: ECS Manager or Editor is nil")
		return
	}

	// Handle player movement input
	is.handlePlayerMovementInput()

	// Handle editor input
	is.handleEditorInput()

	// Handle spawner input
	is.handleSpawnerInput()

	// Handle QuadTree rendering
	is.handleQuadTreeRendering()

	// Handle player play sound input
	is.handlePlayerPlaySound()
}

// handlePlayerMovementInput handles player movement input
func (is *InputSystem) handlePlayerMovementInput() {

	// Get all entities with PlayerComponent
	playerEntityIDs := is.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.PlayerComponent})

	// If there are no entities with PlayerComponent, log an error and return
	if len(playerEntityIDs) == 0 {
		utils.ErrorLogger.Println("InputSystem: No entities with PlayerComponent found")
		return
	}

	// Iterate over all entities with PlayerComponent
	for _, eID := range playerEntityIDs {

		// Get the RigidBodyComponent for the entity
		rbComp, exists := is.manager.GetComponent(eID, ecs.RigidBodyComponent)

		// If the RigidBodyComponent does not exist, log an error and return
		if !exists {
			utils.ErrorLogger.Println("InputSystem: RigidBodyComponent does not exist for entity ", eID)
			return
		}

		// Check if the RigidBodyComponent is of the correct type
		rb, ok := rbComp.(*physicscomponents.RigidBody)

		// If the RigidBodyComponent is not of the correct type, log an error and return
		if !ok {
			utils.ErrorLogger.Println("InputSystem: RigidBodyComponent has wrong type for entity ", eID)
			return
		}

		// Define the movement force
		movementForce := float32(1300.0)

		// Initialize the movement vector
		force := raylib.NewVector2(0, 0)

		// Check if the W key is pressed
		if raylib.IsKeyDown(raylib.KeyW) {
			force.Y = -movementForce
		}

		// Check if the S key is pressed
		if raylib.IsKeyDown(raylib.KeyS) {
			force.Y = movementForce
		}

		// Check if the A key is pressed
		if raylib.IsKeyDown(raylib.KeyA) {
			force.X = -movementForce
		}

		// Check if the D key is pressed
		if raylib.IsKeyDown(raylib.KeyD) {
			force.X = movementForce
		}

		// Apply the movement force to the RigidBody's force
		rb.Force = raylib.Vector2Add(rb.Force, force)
	}
}

// handleEditorInput handles editor input
func (is *InputSystem) handleEditorInput() {

	// Check if the F1 key is pressed
	if raylib.IsKeyPressed(raylib.KeyF1) {

		// Toggle the editor visibility
		is.editor.ToggleVisibility()
	}
}

/*
handleSpawnerInput handles spawner input
creating entities with random positions, velocities, speeds and colors
*/
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

			// Create an entity
			entity := is.manager.CreateEntity()

			// Add components to the entity
			is.manager.AddComponent(entity.ID, ecs.Transform2DComponent, &components.Transform2D{Position: raylib.NewVector2(float32(raylib.GetRandomValue(0, int32(raylib.GetScreenWidth())-1)), float32(raylib.GetRandomValue(0, int32(raylib.GetScreenHeight())-1)))})
			is.manager.AddComponent(entity.ID, ecs.VelocityComponent, &components.Velocity{Vector: raylib.NewVector2(float32(raylib.GetRandomValue(-10, 10)), float32(raylib.GetRandomValue(-10, 10)))})
			is.manager.AddComponent(entity.ID, ecs.SpeedComponent, &components.Speed{Value: float32(raylib.GetRandomValue(50, 200))})
			is.manager.AddComponent(entity.ID, ecs.ColorComponent, &components.Color{Color: color})
		}
	}
}

/*
handleRigidBodySpawner handles rigid body spawner
TODO: Find a way to retrieve the collision system from the Manager
*/

func (is *InputSystem) handleQuadTreeRendering() {
	if raylib.IsKeyPressed(raylib.KeyF2) {

	}
}

/*
handlePlayerPlaySound handles player play sound input
TODO: Remove this function as it's more of a test function
*/
func (is *InputSystem) handlePlayerPlaySound() {
	if raylib.IsKeyPressed(raylib.KeyB) {

		// Get all entities with AudioSourceComponent
		entityIDs := is.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.AudioSourceComponent})

		// Check if there are any entities with AudioSourceComponent
		if len(entityIDs) == 0 {
			utils.InfoLogger.Println("No entities with AudioSourceComponent found")
			return
		}

		// Iterate over all entities with AudioSourceComponent
		for _, eID := range entityIDs {

			// Get the AudioSourceComponent for the entity
			audioComp, exists := is.manager.GetComponent(eID, ecs.AudioSourceComponent)

			// Check if the AudioSourceComponent does not exist, log an error and return
			if !exists {
				utils.InfoLogger.Printf("Entity %d has no AudioSourceComponent", eID)
				return
			}

			// Check if the AudioSourceComponent is of the correct type
			audioSource, ok := audioComp.(*components.AudioSource)

			// Check if the AudioSourceComponent is not of the correct type, log an error and return
			if !ok {
				utils.ErrorLogger.Printf("Entity %d AudioSourceComponent has wrong type", eID)
				return
			}

			// Set the ShouldPlay flag to true
			audioSource.ShouldPlay = true
		}
	}
}

func (is *InputSystem) GetPriority() int {
	return is.priority
}
