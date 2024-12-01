package scenes

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/editor"
	"github.com/webbelito/Fenrir/pkg/physics"
	physicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"
	physicssystems "github.com/webbelito/Fenrir/pkg/physics/systems"
	"github.com/webbelito/Fenrir/pkg/systems"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// GameScene is a struct that represents a game scene
type GameScene struct {
	sceneManager *SceneManager
	manager      *ecs.Manager
	sceneData    *SceneData
	perfMonitor  *editor.PerformanceMonitor

	entities []*ecs.Entity

	// Performance Metrics
	updateDuration time.Duration
	renderDuration time.Duration
	totalDuration  time.Duration

	playerEntity *ecs.Entity
}

// NewGameScene creates a new GameScene
func NewGameScene(sm *SceneManager, m *ecs.Manager, sd *SceneData) *GameScene {
	return &GameScene{
		sceneManager:   sm,
		manager:        m,
		sceneData:      sd,
		updateDuration: 0,
		renderDuration: 0,
		totalDuration:  0,
		entities:       []*ecs.Entity{},
	}
}

// Initialize initializes the Game Scene
func (gs *GameScene) Initialize() {

	// * Editor Init
	editorManager := editor.NewEditorManager(gs.manager, 10)

	// * Resource Manager Init

	// Initialize System specific to GameScene

	// * Input System Init

	inputSystem := systems.NewInputSystem(
		gs.manager,
		editorManager,
		11,
	)

	gs.manager.RegisterSystem(inputSystem, inputSystem.GetPriority())

	// * Movement System Init

	movementSystem := systems.NewMovementSystem(gs.manager, 12)
	gs.manager.RegisterSystem(movementSystem, movementSystem.GetPriority())

	// Add more systems as needed...

	// * Physics Init

	// Initialize the gravity vector (pixels per second i.e 980 pixels per second)
	gravity := raylib.NewVector2(0, 980)

	// Initialize the RigidBodySystem
	rigidBodySystem := physicssystems.NewRigidBodySystem(gs.manager, gravity, 13)
	gs.manager.RegisterSystem(rigidBodySystem, rigidBodySystem.GetPriority())

	// * Collision System Init

	// Initialize the CollisionSystem
	quadBoundary := physics.Rectangle{
		Position: raylib.NewVector2(0, 0),
		Width:    float32(raylib.GetScreenWidth()),
		Height:   float32(raylib.GetScreenHeight()),
	}

	// Initialize the CollisionSystem
	csCapacity := int32(4)
	maxDepth := int32(5)
	capacityDepth := int32(0)

	collisionSystem := physicssystems.NewCollisionSystem(gs.manager, quadBoundary, csCapacity, maxDepth, capacityDepth, 14)
	gs.manager.RegisterSystem(collisionSystem, collisionSystem.GetPriority())

	// * Render Init

	// Initialize the RenderSystem
	screenBoundry := raylib.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(raylib.GetScreenWidth()),
		Height: float32(raylib.GetScreenHeight()),
	}

	renderSystem := systems.NewRenderSystem(gs.manager, screenBoundry, 15)
	gs.manager.RegisterSystem(renderSystem, renderSystem.GetPriority())

	// * Editor Systems

	// Register the Editor Manager
	gs.manager.RegisterSystem(editorManager, editorManager.GetPriority())

	// * Particle System

	// Initialize the Particle System
	particleSystem := systems.NewParticleSystem(gs.manager, 16)
	gs.manager.RegisterSystem(particleSystem, particleSystem.GetPriority())

	// * Particle Render System
	particleRenderSystem := systems.NewParticleRenderSystem(gs.manager, 17)
	gs.manager.RegisterSystem(particleRenderSystem, particleRenderSystem.GetPriority())

	// * Animation System
	animationSystem := systems.NewAnimationSystem(gs.manager, 18)
	gs.manager.RegisterSystem(animationSystem, animationSystem.GetPriority())

	// * Audio System
	audioSystem := systems.NewAudioSystem(gs.manager, 19)
	gs.manager.RegisterSystem(audioSystem, audioSystem.GetPriority())

	// Initialize Entities based on Scene Data
	gs.initializeEntities()

	// Initialize Environment settings
	gs.initializeEnvironment()

	// Spawn 100 entities with random positions and colors
	gs.spawnEntities(25)

}

func (gs *GameScene) Update(dt float64) {

	// Update ECS Manager
	updateStart := time.Now()
	gs.manager.Update(dt)

	// TODO: Remove the temporary input handling
	if raylib.IsKeyPressed(raylib.KeyEscape) {
		err := gs.sceneManager.PushScene("assets/scenes/pause_scene.json")
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}

	// Calculate Performance Metrics
	gs.updateDuration = time.Since(updateStart)
}

func (gs *GameScene) Render() {

	// Render ECS Manager
	renderStart := time.Now()
	gs.manager.Render()

	// Calculate Performance Metrics
	gs.renderDuration = time.Since(renderStart)
	gs.totalDuration = gs.updateDuration + gs.renderDuration

}

// Cleanup cleans up the Game Scene
func (gs *GameScene) Cleanup() {
	// Remove all entities created by this scene
	gs.RemoveAllEntities()

	// TODO: Add the support for cleaning up systems

	// Cleanup resources
	gs.perfMonitor = nil
	raylib.CloseAudioDevice()
}

// Pause pauses the Game Scene
func (gs *GameScene) Pause() {
	// TODO: Implement Pause functionality
	// Pause game logic if necessary
	// For example, stop certain systems or timers
}

// Resume resumes the Game Scene
func (gs *GameScene) Resume() {
	// TODO: Implement Resume functionality
	// Resume game logic if necessary
	// For example, resume certain systems or timers
}

// AddEntity adds an entity to the Game Scene
func (gs *GameScene) AddEntity(e *ecs.Entity) {
	gs.entities = append(gs.entities, e)
}

// RemoveEntity removes an entity from the Game Scene
func (gs *GameScene) RemoveAllEntities() {

	// Iterate over all entities and destroy
	for _, eID := range gs.entities {
		gs.manager.DestroyEntity(eID.ID)
	}

	// Clear the entities slice
	gs.entities = []*ecs.Entity{}
}

// initializeEntities initializes entities based on the scene data
func (gs *GameScene) initializeEntities() {

	// Iterate over all entities in the scene data
	for _, entityData := range gs.sceneData.Entities {

		// Create a new entity
		entity := gs.manager.CreateEntity()

		// Track the entity
		gs.AddEntity(entity)

		// Add components to the entity
		for compName, compData := range entityData.Components {

			switch compName {
			case "Transform2D":

				// Extract the component data
				compMap := compData.(map[string]interface{})

				// Extract the position, rotation and scale data
				position := compMap["position"].(map[string]interface{})
				rotation := compMap["rotation"].(float64)
				scale := compMap["scale"].(map[string]interface{})

				// Create a new Transform2D component
				transform := &components.Transform2D{
					Position: raylib.NewVector2(float32(position["x"].(float64)), float32(position["y"].(float64))),
					Rotation: float32(rotation),
					Scale:    raylib.NewVector2(float32(scale["x"].(float64)), float32(scale["y"].(float64))),
				}

				// Add the Transform2D component to the entity
				gs.manager.AddComponent(entity.ID, ecs.Transform2DComponent, transform)

			case "Sprite":

				// Extract the component data
				compMap := compData.(map[string]interface{})

				// Extract the source rectangle, origin and color data
				sourceRect := compMap["sourceRect"].(map[string]interface{})
				origin := compMap["origin"].(map[string]interface{})
				color := compMap["color"].(string)

				// Create a new Sprite component
				sprite := &components.Sprite{
					TexturePath: compMap["texture_path"].(string),
					SourceRect:  raylib.NewRectangle(float32(sourceRect["x"].(float64)), float32(sourceRect["y"].(float64)), float32(sourceRect["width"].(float64)), float32(sourceRect["height"].(float64))),
					Origin:      raylib.NewVector2(float32(origin["x"].(float64)), float32(origin["y"].(float64))),
					Color:       utils.GetColorFromString(color),
				}

				// Add the Sprite component to the entity
				gs.manager.AddComponent(entity.ID, ecs.SpriteComponent, sprite)

			case "RigidBody":

				// Extract the component data
				compMap := compData.(map[string]interface{})

				// Extract the velocity and acceleration data
				velocity := compMap["velocity"].(map[string]interface{})
				acceleration := compMap["acceleration"].(map[string]interface{})

				// Create a new RigidBody component
				rigidbody := &physicscomponents.RigidBody{
					Mass:         float32(compMap["mass"].(float64)),
					Velocity:     raylib.NewVector2(float32(velocity["x"].(float64)), float32(velocity["y"].(float64))),
					Acceleration: raylib.NewVector2(float32(acceleration["x"].(float64)), float32(acceleration["y"].(float64))),
					Drag:         float32(compMap["drag"].(float64)),
					Restitution:  float32(compMap["restitution"].(float64)),
					IsKinematic:  compMap["is_kinematic"].(bool),
					IsStatic:     compMap["is_static"].(bool),
				}

				// Add the RigidBody component to the entity
				gs.manager.AddComponent(entity.ID, ecs.RigidBodyComponent, rigidbody)

			case "BoxCollider":

				// Extract the component data
				compMap := compData.(map[string]interface{})

				// Extract the size data
				size := compMap["size"].(map[string]interface{})

				// Create a new BoxCollider component
				boxCollider := &physicscomponents.BoxCollider{
					Type: compMap["type"].(string),
					Size: raylib.NewVector2(float32(size["x"].(float64)), float32(size["y"].(float64))),
				}

				// Add the BoxCollider component to the entity
				gs.manager.AddComponent(entity.ID, ecs.BoxColliderComponent, boxCollider)

			case "Color":

				// Extract the component data
				compMap := compData.(map[string]interface{})

				// Create a new Color component
				color := utils.GetColorFromString(compMap["color"].(string))

				// Add the Color component to the entity
				gs.manager.AddComponent(entity.ID, ecs.ColorComponent, &components.Color{Color: color})

			case "Player":

				// Extract the component data
				compMap := compData.(map[string]interface{})

				// Create a new Player component
				player := &components.Player{
					Name: compMap["name"].(string),
				}
				// Add the Player component to the entity
				gs.manager.AddComponent(entity.ID, ecs.PlayerComponent, player)

			case "AudioSource":

				// Extract the component data
				compMap := compData.(map[string]interface{})

				// Extract the file path, volume and isLooping data
				filePath := compMap["file_path"].(string)
				volume := float32(compMap["volume"].(float64))
				isLooping := compMap["is_looping"].(bool)

				// Load the sound file
				sound, err := gs.manager.LoadSound(filePath)
				if err != nil {
					utils.ErrorLogger.Printf("Failed to load sound: %s", err)
					continue
				}

				// Create a new AudioSource component
				audioSource := &components.AudioSource{
					FilePath:  filePath,
					Volume:    volume,
					IsLooping: isLooping,
					Sound:     sound,
				}

				// Add the AudioSource component to the entity
				gs.manager.AddComponent(entity.ID, ecs.AudioSourceComponent, audioSource)

			// Add more components as needed...

			default:
				utils.ErrorLogger.Printf("Component %s not recognized", compName)
			}

			// Assign the player entity
			gs.playerEntity = entity
		}

		// TODO: Remove this temporary code
		// * Particle Emitter example
		/*
			// Create dust emitter entity
			particleEmitter := &components.ParticleEmitter{
				Particles:        []*components.Particle{},
				EmitRate:         10,
				ParticleLifetime: time.Second * 2,
				IsEmitting:       true,
				LastEmitTime:     time.Now(),
			}

			gs.ecsManager.AddComponent(entity.ID, ecs.ParticleEmitterComponent, particleEmitter)

		*/

		// TODO: Remove this temporary code
		//* Animation example

		// Create an Animation component
		frames := []raylib.Rectangle{
			raylib.NewRectangle(0, 0, 32, 32),
			raylib.NewRectangle(32, 0, 32, 32),
			raylib.NewRectangle(64, 0, 32, 32),
			raylib.NewRectangle(96, 0, 32, 32),
		}

		// Create an Animation component
		animation := &components.Animation{
			Frames:        frames,
			CurrentFrame:  0,
			FrameDuration: time.Millisecond * 200,
			IsLooping:     true,
			IsPlaying:     true,
		}

		// Add the Animation component to the entity
		gs.manager.AddComponent(entity.ID, ecs.AnimationComponent, animation)

	}
}

// initializeEnvironment initializes the environment settings
func (gs *GameScene) initializeEnvironment() {

	// Set the background color
	env := gs.sceneData.Environment
	bgColor := utils.GetColorFromString(env.BackgroundColor)
	raylib.ClearBackground(bgColor)

	// Music settings
	if env.Music != "" {
		raylib.InitAudioDevice()
		music := raylib.LoadMusicStream(env.Music)
		raylib.PlayMusicStream(music)

		// Store music reference if needed later for cleanup...
	}
}

// spawnEntities spawns entities with random positions and colors
func (gs *GameScene) spawnEntities(count int) {

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

	// Create entities
	for i := 0; i < count; i++ {

		// Select a random color from the colors slice
		color := colors[raylib.GetRandomValue(0, int32(len(colors)-1))]

		spawnPos := raylib.NewVector2(float32(raylib.GetRandomValue(0, int32(raylib.GetScreenWidth())-1)), float32(raylib.GetRandomValue(0, int32(raylib.GetScreenHeight())-1)))

		// Create an entity with a Transform2D, Rigidbody and Color
		entity := gs.manager.CreateEntity()

		// Add a Transform2D component
		gs.manager.AddComponent(entity.ID, ecs.Transform2DComponent, &components.Transform2D{
			Position: spawnPos,
			Rotation: 0,
			Scale:    raylib.NewVector2(32, 32),
		})

		// Add a Rigidbody component
		gs.manager.AddComponent(entity.ID, ecs.RigidBodyComponent, &physicscomponents.RigidBody{
			Mass:         1.0,
			Velocity:     raylib.NewVector2(0, 0),
			Acceleration: raylib.NewVector2(0, 0),
			Force:        raylib.NewVector2(0, 0),
			Restitution:  0.5,
			Drag:         0.1,
			IsKinematic:  false,
			IsStatic:     false,
		})

		// Add a BoxCollider component
		gs.manager.AddComponent(entity.ID, ecs.BoxColliderComponent, &physicscomponents.BoxCollider{
			Type: "Square",
			Size: raylib.NewVector2(32, 32),
		})

		// Add a Color component
		gs.manager.AddComponent(entity.ID, ecs.ColorComponent, &components.Color{Color: color})
	}
}
