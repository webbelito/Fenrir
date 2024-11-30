package systems

import (
	"sort"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// RenderSystem is a system that renders entities
type RenderSystem struct {
	ScreenCullingRect raylib.Rectangle
	Entities          []EntityData
	manager           *ecs.Manager
	cameraSystem      *CameraSystem
	priority          int
}

// EntityData holds data for rendering an entity
type EntityData struct {
	ID       uint64
	Position raylib.Vector2
	Rotation float32
	Scale    raylib.Vector2
	Color    raylib.Color
	Sprite   *components.Sprite
}

// NewRenderSystem creates a new RenderSystem
func NewRenderSystem(m *ecs.Manager, screenBounds raylib.Rectangle, p int) *RenderSystem {
	return &RenderSystem{
		ScreenCullingRect: screenBounds,
		Entities:          []EntityData{},
		manager:           m,
		priority:          p,
	}
}

func (rs *RenderSystem) Render() {

	// Check if the ECS Manager is nil
	if rs.manager == nil {
		utils.ErrorLogger.Println("ECS Manager is nil")
		return
	}

	// Render entities
	rs.RenderEntities()

}

func (rs *RenderSystem) RenderEntities() {
	// Retrieve cameraComp from Manager
	// TODO We are assuming that the camera is the first entity
	cameraComp, exist := rs.manager.GetComponent(0, ecs.CameraComponent)

	// Check if the camera exists
	if !exist {
		utils.ErrorLogger.Println("RenderSystem: Camera not found")
		return
	}

	// Cast the camera component to a Camera component
	camera, ok := cameraComp.(*components.Camera)

	// Check if the cast was successful
	if !ok {
		utils.ErrorLogger.Println("RenderSystem: Invalid Camera component")
		return
	}

	// Create a Camera2D object
	cam := raylib.Camera2D{
		Offset:   camera.Offset,
		Target:   camera.Target,
		Rotation: 0,
		Zoom:     camera.Zoom,
	}

	// Begin 2D mode with the camera
	raylib.BeginMode2D(cam)

	// Retrieve entities with Transform2D and Color components
	entityIDs := rs.manager.GetEntitiesWithComponents([]ecs.ComponentType{
		ecs.Transform2DComponent,
		ecs.ColorComponent,
	})

	// Temporary slice to hold entities for rendering
	var entities []EntityData

	for _, eID := range entityIDs {
		// Retrieve Transform2D component
		transformComp, exists := rs.manager.GetComponent(eID, ecs.Transform2DComponent)
		if !exists {
			continue
		}
		transform, ok := transformComp.(*components.Transform2D)
		if !ok {
			utils.ErrorLogger.Printf("RenderSystem: Entity %d has invalid Transform2D component\n", eID)
			continue
		}

		// Check if the entity is within the screen bounds
		if !raylib.CheckCollisionPointRec(transform.Position, rs.ScreenCullingRect) {
			continue
		}

		// Retrieve Color component
		colorComp, exists := rs.manager.GetComponent(eID, ecs.ColorComponent)
		if !exists {
			continue
		}
		color, ok := colorComp.(*components.Color)
		if !ok {
			utils.ErrorLogger.Printf("RenderSystem: Entity %d has invalid Color component\n", eID)
			continue
		}

		// Retrieve Sprite component, if any
		spriteComp, spriteExists := rs.manager.GetComponent(eID, ecs.SpriteComponent)
		var sprite *components.Sprite
		if spriteExists {
			sprite, ok = spriteComp.(*components.Sprite)
			if !ok {
				utils.ErrorLogger.Printf("RenderSystem: Entity %d has invalid Sprite component\n", eID)
				continue
			}
		}

		// Append entity data for rendering
		entities = append(entities, EntityData{
			ID:       eID,
			Position: transform.Position,
			Rotation: transform.Rotation,
			Scale:    transform.Scale,
			Color:    color.Color,
			Sprite:   sprite,
		})
	}

	// Sort entities by ID to ensure consistent rendering order
	sort.SliceStable(entities, func(i, j int) bool {
		return entities[i].ID < entities[j].ID
	})

	// Render entities
	for _, entity := range entities {
		if entity.Sprite == nil {
			// Render a rectangle if the sprite is nil
			raylib.DrawRectangleV(entity.Position, entity.Scale, entity.Color)
			continue
		}

		// Retrieve the texture from the Resources Manager
		texture, texExists := rs.manager.GetTexture(entity.Sprite.TexturePath)
		if !texExists {
			// Attempt to load the texture if not already loaded
			var err error
			texture, err = rs.manager.LoadTexture(entity.Sprite.TexturePath)
			if err != nil {
				utils.ErrorLogger.Printf("RenderSystem: Failed to load texture: %s\n", entity.Sprite.TexturePath)
				continue
			}
		}

		// Define source and destination rectangles
		sourceRect := entity.Sprite.SourceRect
		destRect := raylib.Rectangle{
			X:      entity.Position.X,
			Y:      entity.Position.Y,
			Width:  entity.Scale.X,
			Height: entity.Scale.Y,
		}

		// Draw the sprite
		raylib.DrawTexturePro(
			texture,
			sourceRect,
			destRect,
			entity.Sprite.Origin,
			entity.Rotation,
			entity.Color,
		)
	}

	// End 2D mode
	raylib.EndMode2D()
}

/*
GetPriority returns the priority of the system
*/
func (rs *RenderSystem) GetPriority() int {
	return rs.priority
}

/*
SetCameraSystem sets the CameraSystem for the RenderSystem
*/
func (rs *RenderSystem) SetCameraSystem(cs *CameraSystem) {
	rs.cameraSystem = cs
}
