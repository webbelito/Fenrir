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
	sceneManager *SceneManager
	ecsManager   *ecs.ECSManager
	sceneData    *SceneData
	perfMonitor  *editor.PerformanceMonitor

	entities      []*ecs.Entity
	logicSystems  []systeminterfaces.Updatable
	renderSystems []systeminterfaces.Renderable

	// Performance Metrics
	updateDuration time.Duration
	renderDuration time.Duration
	totalDuration  time.Duration
}

func NewGameScene(sm *SceneManager, em *ecs.ECSManager, sd *SceneData) *GameScene {
	return &GameScene{
		sceneManager:   sm,
		ecsManager:     em,
		sceneData:      sd,
		updateDuration: 0,
		renderDuration: 0,
		totalDuration:  0,
		entities:       []*ecs.Entity{},
		logicSystems:   []systeminterfaces.Updatable{},
		renderSystems:  []systeminterfaces.Renderable{},
	}
}

func (gs *GameScene) Init() {

	// Initialize the Editor
	gameEditor := editor.NewEditor(gs.ecsManager)

	// Initialize the Resource Manager
	resourceManager := resources.NewResourceManager()

	// Initialize Performance Monitor
	gs.perfMonitor = editor.NewPerformanceMonitor(raylib.NewVector2(1600, 10))

	// Initialize System specific to GameScene
	inputSystem := systems.NewInputSystem(
		gs.ecsManager,
		gameEditor,
	)

	gs.ecsManager.AddLogicSystem(inputSystem, 0)
	gs.logicSystems = append(gs.logicSystems, inputSystem)

	movementSystem := systems.NewMovementSystem(gs.ecsManager)
	gs.ecsManager.AddLogicSystem(movementSystem, 1)
	gs.logicSystems = append(gs.logicSystems, movementSystem)

	// Add more systems as needed...

	// Initialize the gravity vector (pixels per second i.e 980 pixels per second)
	gravity := raylib.NewVector2(0, 980)

	// Initialize the RigidBodySystem
	rigidBodySystem := physicssystems.NewRigidBodySystem(gs.ecsManager, gravity)
	gs.ecsManager.AddLogicSystem(rigidBodySystem, 2)
	gs.logicSystems = append(gs.logicSystems, rigidBodySystem)

	// Initialize the CollisionSystem
	quadBoundary := physics.Rectangle{
		Position: raylib.NewVector2(0, 0),
		Width:    float32(raylib.GetScreenWidth()),
		Height:   float32(raylib.GetScreenHeight()),
	}

	csCapacity := int32(4)
	maxDepth := int32(5)
	capacityDepth := int32(0)

	collisionSystem := physicssystems.NewCollisionSystem(gs.ecsManager, quadBoundary, csCapacity, maxDepth, capacityDepth)
	gs.ecsManager.AddLogicSystem(collisionSystem, 3)
	gs.logicSystems = append(gs.logicSystems, collisionSystem)

	renderSystem := systems.NewRenderSystem(gs.ecsManager, raylib.NewRectangle(0, 0, float32(raylib.GetScreenWidth()), float32(raylib.GetScreenHeight())), resourceManager)
	gs.ecsManager.AddRenderSystem(renderSystem, 0)
	gs.renderSystems = append(gs.renderSystems, renderSystem)

	// Initialize Entities based on Scene Data
	gs.initializeEntities()

	// Initialize Environment settings
	gs.initializeEnvironment()

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

	// Update Performance Monitor
	perfData := &editor.PerformanceMonitorData{
		FPS:            raylib.GetFPS(),
		UpdateDuration: gs.updateDuration,
		RenderDuration: gs.renderDuration,
		TotalDuration:  gs.totalDuration,
	}

	// Update Performance Monitor
	gs.perfMonitor.Update(perfData)

	// Render Performance Monitor
	gs.perfMonitor.Draw(perfData)

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

			// Add more components as needed...

			default:
				utils.ErrorLogger.Printf("Component %s not recognized", compName)
			}
		}
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
