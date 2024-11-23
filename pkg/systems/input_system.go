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
	ecsManager        *ecs.ECSManager
	entitiesManager   *ecs.EntitiesManager
	componentsManager *ecs.ComponentsManager
	editorManager     *editor.EditorManager
	priority          int
}

func NewInputSystem(ecsM *ecs.ECSManager, e *editor.EditorManager, p int) *InputSystem {
	return &InputSystem{
		ecsManager:        ecsM,
		entitiesManager:   ecsM.GetEntitiesManager(),
		componentsManager: ecsM.GetComponentsManager(),
		editorManager:     e,
		priority:          p,
	}
}

func (is *InputSystem) Update(dt float64) {

	if is.ecsManager == nil || is.entitiesManager == nil || is.componentsManager == nil || is.editorManager == nil {
		utils.ErrorLogger.Println("InputSystem: EntitiesManager, ComponentsManager or EditorManager is nil")
		return
	}

	// Handle player movement input
	is.handlePlayerMovementInput()

	// Handle editor input
	is.handleEditorInput()

	// Handle spawner input
	is.handleSpawnerInput()

	// Handle rigid body spawner
	is.handleRigidBodySpawner()

	// Handle QuadTree rendering
	is.handleQuadTreeRendering()

	// TODO: Change this to something proper
	is.handlePlayerPlaySound()
}

func (is *InputSystem) handlePlayerMovementInput() {

	playerComps, playerCompsExists := is.componentsManager.Components[ecs.PlayerComponent]

	if !playerCompsExists {
		utils.ErrorLogger.Println("InputSystem: No player components found")
		return
	}

	for entity := range playerComps {
		rbComponent, rbCompExists := is.ecsManager.GetComponent(entity, ecs.RigidBodyComponent)

		if !rbCompExists {
			continue
		}

		rb, rbExists := rbComponent.(*physicscomponents.RigidBody)

		if !rbExists {
			continue
		}

		// Define the movement force
		movementForce := float32(1300.0)

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
		utils.InfoLogger.Println("InputSystem: Toggling editor visibility")
		is.editorManager.ToggleVisibility()
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
			entity := is.ecsManager.CreateEntity()
			is.ecsManager.AddComponent(entity.ID, ecs.Transform2DComponent, &components.Transform2D{Position: raylib.NewVector2(float32(raylib.GetRandomValue(0, int32(raylib.GetScreenWidth())-1)), float32(raylib.GetRandomValue(0, int32(raylib.GetScreenHeight())-1)))})
			is.ecsManager.AddComponent(entity.ID, ecs.VelocityComponent, &components.Velocity{Vector: raylib.NewVector2(float32(raylib.GetRandomValue(-10, 10)), float32(raylib.GetRandomValue(-10, 10)))})
			is.ecsManager.AddComponent(entity.ID, ecs.SpeedComponent, &components.Speed{Value: float32(raylib.GetRandomValue(50, 200))})
			is.ecsManager.AddComponent(entity.ID, ecs.ColorComponent, &components.Color{Color: color})
		}
	}
}

func (is *InputSystem) handleRigidBodySpawner() {

	if raylib.IsKeyPressed(raylib.KeyR) {

		// Create a rigid body entity
		rigidBodyEntity := is.ecsManager.CreateEntity()

		// Add a rigid body component to the rigid body entity
		is.ecsManager.AddComponent(rigidBodyEntity.ID, ecs.RigidBodyComponent, &physicscomponents.RigidBody{
			Mass:         1,
			Velocity:     raylib.NewVector2(0, 0),
			Acceleration: raylib.NewVector2(0, 0),
			Force:        raylib.NewVector2(0, 0),
			IsKinematic:  false,
			IsStatic:     false,
		})

	}
}

func (is *InputSystem) handleQuadTreeRendering() {
	if raylib.IsKeyPressed(raylib.KeyF2) {
		// TODO: Find a way to retrieve the collision system from the ECSManager
	}
}

func (is *InputSystem) handlePlayerPlaySound() {
	if raylib.IsKeyPressed(raylib.KeyB) {

		playerEntities := is.ecsManager.GetComponentsManager().GetEntitiesWithComponents([]ecs.ComponentType{ecs.AudioSourceComponent})
		for _, entityID := range playerEntities {
			audioComp, audioCompExists := is.ecsManager.GetComponent(entityID, ecs.AudioSourceComponent)
			if !audioCompExists {
				utils.InfoLogger.Printf("Entity %d has no AudioSourceComponent", entityID)
				continue
			}

			audioSource, ok := audioComp.(*components.AudioSource)
			if !ok {
				utils.ErrorLogger.Printf("Entity %d AudioSourceComponent has wrong type", entityID)
				continue
			}

			audioSource.ShouldPlay = true
		}
	}
}

func (is *InputSystem) GetPriority() int {
	return is.priority
}
