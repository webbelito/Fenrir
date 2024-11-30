package scenes

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/editor"
	systeminterfaces "github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/physics"
	physicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"
	physicssystems "github.com/webbelito/Fenrir/pkg/physics/systems"
	"github.com/webbelito/Fenrir/pkg/resources"
	"github.com/webbelito/Fenrir/pkg/systems"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type GameScene struct {
	sceneManager    *SceneManager
	resourceManager *resources.ResourcesManager
	ecsManager      *ecs.ECSManager
	sceneData       *SceneData
	perfMonitor     *editor.PerformanceMonitor

	entities      []*ecs.Entity
	logicSystems  []systeminterfaces.UpdatableSystemInterface
	renderSystems []systeminterfaces.RenderableSystemInterface

	// Performance Metrics
	updateDuration time.Duration
	renderDuration time.Duration
	totalDuration  time.Duration

	playerEntity *ecs.Entity
}

func NewGameScene(sm *SceneManager, em *ecs.ECSManager, sd *SceneData) *GameScene {
	return &GameScene{
		sceneManager:    sm,
		resourceManager: resources.NewResourceManager(),
		ecsManager:      em,
		sceneData:       sd,
		updateDuration:  0,
		renderDuration:  0,
		totalDuration:   0,
		entities:        []*ecs.Entity{},
		logicSystems:    []systeminterfaces.UpdatableSystemInterface{},
		renderSystems:   []systeminterfaces.RenderableSystemInterface{},
	}
}

func (gs *GameScene) Initialize() {

	// * Editor Init
	editorManager := editor.NewEditorManager(gs.ecsManager, 1)

	// * Resource Manager Init

	// Initialize System specific to GameScene

	// * Input System Init

	inputSystem := systems.NewInputSystem(
		gs.ecsManager,
		editorManager,
		2,
	)

	gs.ecsManager.AddLogicSystem(inputSystem, inputSystem.GetPriority())
	gs.logicSystems = append(gs.logicSystems, inputSystem)

	// * Movement System Init

	movementSystem := systems.NewMovementSystem(gs.ecsManager, 3)
	gs.ecsManager.AddLogicSystem(movementSystem, movementSystem.GetPriority())
	gs.logicSystems = append(gs.logicSystems, movementSystem)

	// Add more systems as needed...

	// * Physics Init

	// Initialize the gravity vector (pixels per second i.e 980 pixels per second)
	gravity := raylib.NewVector2(0, 980)

	// Initialize the RigidBodySystem
	rigidBodySystem := physicssystems.NewRigidBodySystem(gs.ecsManager, gravity, 4)
	gs.ecsManager.AddLogicSystem(rigidBodySystem, rigidBodySystem.GetPriority())
	gs.logicSystems = append(gs.logicSystems, rigidBodySystem)

	// * Collision System Init

	// Initialize the CollisionSystem
	quadBoundary := physics.Rectangle{
		Position: raylib.NewVector2(0, 0),
		Width:    float32(raylib.GetScreenWidth()),
		Height:   float32(raylib.GetScreenHeight()),
	}

	csCapacity := int32(4)
	maxDepth := int32(5)
	capacityDepth := int32(0)

	collisionSystem := physicssystems.NewCollisionSystem(gs.ecsManager, quadBoundary, csCapacity, maxDepth, capacityDepth, 5)
	gs.ecsManager.AddLogicSystem(collisionSystem, collisionSystem.GetPriority())
	gs.logicSystems = append(gs.logicSystems, collisionSystem)

	// * Render Init

	screenBoundry := raylib.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(raylib.GetScreenWidth()),
		Height: float32(raylib.GetScreenHeight()),
	}

	renderSystem := systems.NewRenderSystem(gs.ecsManager, screenBoundry, gs.resourceManager, 0)
	gs.ecsManager.AddRenderSystem(renderSystem, renderSystem.GetPriority())
	gs.renderSystems = append(gs.renderSystems, renderSystem)

	// * Editor Systems

	gs.ecsManager.AddLogicSystem(editorManager, editorManager.GetPriority())
	gs.logicSystems = append(gs.logicSystems, editorManager)

	gs.ecsManager.AddRenderSystem(editorManager, 1)
	gs.renderSystems = append(gs.renderSystems, editorManager)

	// * Particle System

	particleSystem := systems.NewParticleSystem(gs.ecsManager, 6)
	gs.ecsManager.AddLogicSystem(particleSystem, particleSystem.GetPriority())
	gs.logicSystems = append(gs.logicSystems, particleSystem)

	particleRenderSystem := systems.NewParticleRenderSystem(gs.ecsManager, 2)
	gs.ecsManager.AddRenderSystem(particleRenderSystem, particleRenderSystem.GetPriority())
	gs.renderSystems = append(gs.renderSystems, particleRenderSystem)

	// * Animation System
	animationSystem := systems.NewAnimationSystem(gs.ecsManager, 7)
	gs.ecsManager.AddLogicSystem(animationSystem, animationSystem.GetPriority())
	gs.logicSystems = append(gs.logicSystems, animationSystem)

	// * Camera System
	// TODO: Move this to a persistent system
	cameraSystem := systems.NewCameraSystem(gs.ecsManager, 8)
	gs.ecsManager.AddLogicSystem(cameraSystem, cameraSystem.GetPriority())
	gs.logicSystems = append(gs.logicSystems, cameraSystem)

	// Assign the camera system to the Render System
	renderSystem.SetCameraSystem(cameraSystem)

	// * Audio System
	audioSystem := systems.NewAudioSystem(gs.ecsManager, gs.resourceManager, 9)
	gs.ecsManager.AddLogicSystem(audioSystem, audioSystem.GetPriority())
	gs.logicSystems = append(gs.logicSystems, audioSystem)

	// Initialize Entities based on Scene Data
	gs.initializeEntities()

	// Initialize Environment settings
	gs.initializeEnvironment()

	// Spawn 100 entities with random positions and colors
	gs.spawnEntities(25)

	// Assign the player entity as the camera's owner
	cameraSystem.SetOwner(gs.playerEntity.ID)
}

func (gs *GameScene) Update(dt float64) {

	// Update ECS Manager
	updateStart := time.Now()
	gs.ecsManager.UpdateLogicSystems(dt)

	// TODO: Remove the temporary input handling
	if raylib.IsKeyPressed(raylib.KeyEscape) {
		err := gs.sceneManager.PushScene("assets/scenes/pause_scene.json")
		if err != nil {
			utils.ErrorLogger.Println("Failed to change scene: ", err)
		}
	}

	gs.updateDuration = time.Since(updateStart)
}

func (gs *GameScene) Render() {

	// Render ECS Manager
	renderStart := time.Now()
	gs.ecsManager.UpdateRenderSystems()

	// Calculate Performance Metrics
	gs.renderDuration = time.Since(renderStart)
	gs.totalDuration = gs.updateDuration + gs.renderDuration

}

func (gs *GameScene) Cleanup() {
	// Remove all entities created by this scene
	gs.RemoveAllEntities()

	// Remove and cleanup logic system
	for _, system := range gs.logicSystems {
		gs.ecsManager.RemoveLogicSystem(system)
	}
	gs.logicSystems = nil

	// Remove and cleanup render system
	for _, system := range gs.renderSystems {
		gs.ecsManager.RemoveRenderSystem(system)
	}

	// Cleanup resources
	gs.perfMonitor = nil
	raylib.CloseAudioDevice()
}

func (gs *GameScene) Pause() {
	// TODO: Implement Pause functionality
	// Pause game logic if necessary
	// For example, stop certain systems or timers
}

func (gs *GameScene) Resume() {
	// TODO: Implement Resume functionality
	// Resume game logic if necessary
	// For example, resume certain systems or timers
}

func (gs *GameScene) AddEntity(e *ecs.Entity) {
	gs.entities = append(gs.entities, e)
}

func (gs *GameScene) RemoveAllEntities() {
	for _, eID := range gs.entities {
		gs.ecsManager.DestroyEntity(eID.ID)
	}

	gs.entities = []*ecs.Entity{}
}

func (gs *GameScene) initializeEntities() {

	// Iterate over all entities in the scene data
	for _, entityData := range gs.sceneData.Entities {

		// Create a new entity
		entity := gs.ecsManager.CreateEntity()

		// Track the entity
		gs.AddEntity(entity)

		// Add components to the entity
		for compName, compData := range entityData.Components {

			switch compName {
			case "Transform2D":
				compMap := compData.(map[string]interface{})
				position := compMap["position"].(map[string]interface{})
				rotation := compMap["rotation"].(float64)
				scale := compMap["scale"].(map[string]interface{})

				transform := &components.Transform2D{
					Position: raylib.NewVector2(float32(position["x"].(float64)), float32(position["y"].(float64))),
					Rotation: float32(rotation),
					Scale:    raylib.NewVector2(float32(scale["x"].(float64)), float32(scale["y"].(float64))),
				}

				gs.ecsManager.AddComponent(entity.ID, ecs.Transform2DComponent, transform)

			case "Sprite":
				compMap := compData.(map[string]interface{})
				sourceRect := compMap["sourceRect"].(map[string]interface{})
				origin := compMap["origin"].(map[string]interface{})
				color := compMap["color"].(string)

				sprite := &components.Sprite{
					TexturePath: compMap["texture_path"].(string),
					SourceRect:  raylib.NewRectangle(float32(sourceRect["x"].(float64)), float32(sourceRect["y"].(float64)), float32(sourceRect["width"].(float64)), float32(sourceRect["height"].(float64))),
					Origin:      raylib.NewVector2(float32(origin["x"].(float64)), float32(origin["y"].(float64))),
					Color:       utils.GetColorFromString(color),
				}

				gs.ecsManager.AddComponent(entity.ID, ecs.SpriteComponent, sprite)

			case "RigidBody":
				compMap := compData.(map[string]interface{})
				velocity := compMap["velocity"].(map[string]interface{})
				acceleration := compMap["acceleration"].(map[string]interface{})

				rigidbody := &physicscomponents.RigidBody{
					Mass:         float32(compMap["mass"].(float64)),
					Velocity:     raylib.NewVector2(float32(velocity["x"].(float64)), float32(velocity["y"].(float64))),
					Acceleration: raylib.NewVector2(float32(acceleration["x"].(float64)), float32(acceleration["y"].(float64))),
					Drag:         float32(compMap["drag"].(float64)),
					Restitution:  float32(compMap["restitution"].(float64)),
					IsKinematic:  compMap["is_kinematic"].(bool),
					IsStatic:     compMap["is_static"].(bool),
				}

				gs.ecsManager.AddComponent(entity.ID, ecs.RigidBodyComponent, rigidbody)

			case "BoxCollider":
				compMap := compData.(map[string]interface{})
				size := compMap["size"].(map[string]interface{})
				boxCollider := &physicscomponents.BoxCollider{
					Type: compMap["type"].(string),
					Size: raylib.NewVector2(float32(size["x"].(float64)), float32(size["y"].(float64))),
				}

				gs.ecsManager.AddComponent(entity.ID, ecs.BoxColliderComponent, boxCollider)

			case "Color":
				compMap := compData.(map[string]interface{})
				color := utils.GetColorFromString(compMap["color"].(string))
				gs.ecsManager.AddComponent(entity.ID, ecs.ColorComponent, &components.Color{Color: color})

			case "Player":
				compMap := compData.(map[string]interface{})
				player := &components.Player{
					Name: compMap["name"].(string),
				}
				gs.ecsManager.AddComponent(entity.ID, ecs.PlayerComponent, player)
			case "AudioSource":
				compMap := compData.(map[string]interface{})
				filePath := compMap["file_path"].(string)
				volume := float32(compMap["volume"].(float64))
				isLooping := compMap["is_looping"].(bool)

				sound, err := gs.resourceManager.LoadSound(filePath)
				if err != nil {
					utils.ErrorLogger.Printf("Failed to load sound: %s", err)
					continue
				}

				audioSource := &components.AudioSource{
					FilePath:  filePath,
					Volume:    volume,
					IsLooping: isLooping,
					Sound:     sound,
				}

				gs.ecsManager.AddComponent(entity.ID, ecs.AudioSourceComponent, audioSource)

			// Add more components as needed...

			default:
				utils.ErrorLogger.Printf("Component %s not recognized", compName)
			}

			gs.playerEntity = entity
		}

		// TODO: Remove this temporary code

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

		frames := []raylib.Rectangle{
			raylib.NewRectangle(0, 0, 32, 32),
			raylib.NewRectangle(32, 0, 32, 32),
			raylib.NewRectangle(64, 0, 32, 32),
			raylib.NewRectangle(96, 0, 32, 32),
		}

		animation := &components.Animation{
			Frames:        frames,
			CurrentFrame:  0,
			FrameDuration: time.Millisecond * 200,
			IsLooping:     true,
			IsPlaying:     true,
		}

		gs.ecsManager.AddComponent(entity.ID, ecs.AnimationComponent, animation)

	}
}

func (gs *GameScene) initializeEnvironment() {
	env := gs.sceneData.Environment
	bgColor := utils.GetColorFromString(env.BackgroundColor)
	raylib.ClearBackground(bgColor)

	if env.Music != "" {
		raylib.InitAudioDevice()
		music := raylib.LoadMusicStream(env.Music)
		raylib.PlayMusicStream(music)

		// Store music reference if needed later for cleanup...
	}
}

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
		entity := gs.ecsManager.CreateEntity()

		// Add a Transform2D component
		gs.ecsManager.AddComponent(entity.ID, ecs.Transform2DComponent, &components.Transform2D{
			Position: spawnPos,
			Rotation: 0,
			Scale:    raylib.NewVector2(32, 32),
		})

		// Add a Rigidbody component
		gs.ecsManager.AddComponent(entity.ID, ecs.RigidBodyComponent, &physicscomponents.RigidBody{
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
		gs.ecsManager.AddComponent(entity.ID, ecs.BoxColliderComponent, &physicscomponents.BoxCollider{
			Type: "Square",
			Size: raylib.NewVector2(32, 32),
		})

		// Add a Color component
		gs.ecsManager.AddComponent(entity.ID, ecs.ColorComponent, &components.Color{Color: color})
	}

}
