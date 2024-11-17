package main

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"

	"github.com/webbelito/Fenrir/pkg/editor"

	"github.com/webbelito/Fenrir/pkg/systems"

	phsyics "github.com/webbelito/Fenrir/pkg/physics"
	phsyicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"
	phsyicssystems "github.com/webbelito/Fenrir/pkg/physics/systems"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

func main() {

	// Initialize Raylib
	raylib.InitWindow(1920, 1280, "Fenrir Test")
	defer raylib.CloseWindow()

	raylib.SetTargetFPS(60)

	// Initialize ECS Manager
	ecsManager := ecs.NewECSManager()

	// Initialize the Editor
	gameEditor := editor.NewEditor(ecsManager)

	// Add systems to the ECS Manager
	ecsManager.AddSystem(&systems.InputSystem{
		Editor:     gameEditor,
		EcsManager: ecsManager,
	}, 0)

	ecsManager.AddSystem(&systems.MovementSystem{}, 1)

	// TODO: Implement a camera and follow the position of the camera with the screen culling rect
	ecsManager.AddSystem(&systems.RenderSystem{ScreenCullingRect: raylib.Rectangle{
		X:      0,
		Y:      0,
		Width:  float32(raylib.GetScreenWidth()),
		Height: float32(raylib.GetScreenHeight()),
	}}, 2)

	// Initialize the gravity vector (pixels per second i.e 980 pixels per second)
	gravity := raylib.NewVector2(0, 980)

	// Initialize the RigidBodySystem
	rigidBodySystem := phsyicssystems.NewRigidBodySystem(gravity)
	ecsManager.AddSystem(rigidBodySystem, 3)

	// Initialize the CollisionSystem
	quadBoundary := phsyics.Rectangle{
		Position: raylib.NewVector2(0, 0),
		Width:    float32(raylib.GetScreenWidth()),
		Height:   float32(raylib.GetScreenHeight()),
	}

	csCapacity := int32(4)
	maxDepth := int32(5)
	capacityDepth := int32(0)

	collisionSystem := phsyicssystems.NewCollisionSystem(quadBoundary, csCapacity, maxDepth, capacityDepth)
	ecsManager.AddSystem(collisionSystem, 4)

	// Create a player entity
	player := ecsManager.CreateEntity()

	ecsManager.AddComponent(player.ID, ecs.Transform2DComponent, &components.Transform2D{
		Position: raylib.NewVector2(100, 100),
		Rotation: 0,
		Scale:    raylib.NewVector2(15, 15),
	})
	ecsManager.AddComponent(player.ID, ecs.ColorComponent, &components.Color{Color: raylib.Red})
	ecsManager.AddComponent(player.ID, ecs.PlayerComponent, &components.Player{Name: "Webbelito"})

	// Add RigidBody component to the player entity
	ecsManager.AddComponent(player.ID, ecs.RigidBodyComponent, &phsyicscomponents.RigidBody{
		Mass:         1,
		Velocity:     raylib.NewVector2(0, 0),
		Acceleration: raylib.NewVector2(0, 0),
		Force:        raylib.NewVector2(0, 0),
		Drag:         1,
		Restitution:  0.5,
		IsKinematic:  false,
		IsStatic:     false,
	})
	ecsManager.AddComponent(player.ID, ecs.BoxColliderComponent, &phsyicscomponents.BoxCollider{
		Type: "Square",
		Size: raylib.NewVector2(16, 16),
	})

	for i := 0; i < 100; i++ {

		// Create a random rigid body entity
		rigidBodyEntity := ecsManager.CreateEntity()
		ecsManager.AddComponent(rigidBodyEntity.ID, ecs.Transform2DComponent, &components.Transform2D{
			Position: raylib.NewVector2(
				float32(raylib.GetRandomValue(100, 1000)),
				float32(raylib.GetRandomValue(100, 1000)),
			),
			Rotation: 0,
			Scale:    raylib.NewVector2(15, 15),
		})

		ecsManager.AddComponent(rigidBodyEntity.ID, ecs.ColorComponent, &components.Color{Color: raylib.Blue})
		ecsManager.AddComponent(rigidBodyEntity.ID, ecs.RigidBodyComponent, &phsyicscomponents.RigidBody{
			Mass:         0.01,
			Velocity:     raylib.NewVector2(0, 0),
			Acceleration: raylib.NewVector2(0, 0),
			Force:        raylib.NewVector2(0, 0),
			Drag:         0.1,
			Restitution:  0.5,
			IsKinematic:  false,
			IsStatic:     false,
		})

		ecsManager.AddComponent(rigidBodyEntity.ID, ecs.BoxColliderComponent, &phsyicscomponents.BoxCollider{
			Type: "Square",
			Size: raylib.NewVector2(15, 15),
		})

	}
	// Main game loop
	for !raylib.WindowShouldClose() {

		// Get the time taken for the last frame
		deltaTime := raylib.GetFrameTime()

		// Update ECS entities
		updateStart := time.Now()
		ecsManager.Update(float64(deltaTime))
		updateDuration := time.Since(updateStart)

		// Begin drawing
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		// Render ECS entities
		renderStart := time.Now()
		ecsManager.Render()
		renderDuration := time.Since(renderStart)

		// Calculate the total time taken for the update and render steps
		totalDuration := time.Since(updateStart)

		// Render Editor Overlay
		gameEditor.Draw(&editor.PerformanceMonitorData{FPS: raylib.GetFPS(), UpdateDuration: updateDuration, RenderDuration: renderDuration, TotalDuration: totalDuration})

		// End drawing
		raylib.EndDrawing()

	}
}
