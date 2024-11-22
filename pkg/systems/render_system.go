package systems

import (
	"sort"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/resources"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type RenderSystem struct {
	ScreenCullingRect raylib.Rectangle
	Entities          []EntityData
	ecsManager        *ecs.ECSManager
	entitiesManager   *ecs.EntitiesManager
	componentsManager *ecs.ComponentsManager
	resourcesManager  *resources.ResourcesManager
	cameraSystem      *CameraSystem
	priority          int
}

type EntityData struct {
	ID       uint64
	Position raylib.Vector2
	Rotation float32
	Scale    raylib.Vector2
	Color    raylib.Color
	Sprite   *components.Sprite
}

func NewRenderSystem(ecsM *ecs.ECSManager, screenBounds raylib.Rectangle, rm *resources.ResourcesManager, p int) *RenderSystem {
	return &RenderSystem{
		ScreenCullingRect: screenBounds,
		Entities:          []EntityData{},
		ecsManager:        ecsM,
		entitiesManager:   ecsM.GetEntitiesManager(),
		componentsManager: ecsM.GetComponentsManager(),
		resourcesManager:  rm,
		priority:          p,
	}
}

func (rs *RenderSystem) Render() {
	if rs.ecsManager == nil || rs.entitiesManager == nil || rs.componentsManager == nil || rs.resourcesManager == nil {
		utils.ErrorLogger.Println("RenderSystem: ECSManager or EntitiesManager or ComponentsManager or ResourcesManager is nil")
		return
	}

	rs.RenderEntities()

}

func (rs *RenderSystem) RenderEntities() {

	camera := rs.cameraSystem.GetCamera()

	cam := raylib.Camera2D{
		Offset:   rs.cameraSystem.camera.Offset,
		Target:   rs.cameraSystem.camera.Target,
		Rotation: 0,
		Zoom:     camera.Zoom,
	}

	raylib.BeginMode2D(cam)

	transformComps, transformCompsExists := rs.componentsManager.Components[ecs.Transform2DComponent]
	colorComps, colorCompsExists := rs.componentsManager.Components[ecs.ColorComponent]
	spriteComps := rs.componentsManager.Components[ecs.SpriteComponent]

	if !transformCompsExists || !colorCompsExists {
		return
	}

	// Clear the entities slice while retaining capacity
	rs.Entities = rs.Entities[:0]

	// Collect entities that have Position , Color and Sprite components
	entities := rs.componentsManager.GetEntitiesWithComponents([]ecs.ComponentType{
		ecs.Transform2DComponent,
		ecs.ColorComponent,
	})

	for _, entity := range entities {
		transform, _ := transformComps[entity].(*components.Transform2D)
		colorComp, _ := colorComps[entity].(*components.Color)

		// Check if the entity has a Sprite component
		spriteComp, spriteExists := spriteComps[entity].(*components.Sprite)

		// Check if the entity is within the screen bounds
		if !raylib.CheckCollisionPointRec(transform.Position, rs.ScreenCullingRect) {
			continue
		}

		// Add the entity to the slice if it does not have a sprite component
		if !spriteExists {
			rs.Entities = append(rs.Entities, EntityData{
				ID:       entity,
				Position: transform.Position,
				Rotation: transform.Rotation,
				Scale:    transform.Scale,
				Color:    colorComp.Color,
				Sprite:   nil,
			})
		}

		// Add the entity to the slice
		rs.Entities = append(rs.Entities, EntityData{
			ID:       entity,
			Position: transform.Position,
			Rotation: transform.Rotation,
			Scale:    transform.Scale,
			Color:    colorComp.Color,
			Sprite:   spriteComp,
		})
	}

	// Sort entities by ID to ensure consistent rendering order
	sort.SliceStable(rs.Entities, func(i int, j int) bool {
		return rs.Entities[i].ID < rs.Entities[j].ID
	})

	// Render entities
	for _, entity := range rs.Entities {
		if entity.Sprite == nil {
			// Fallback to rendering a rectangle if the sprite is nil
			raylib.DrawRectangleV(entity.Position, entity.Scale, entity.Color)
			continue
		}

		// Retrieve the texture from the Resources Manager
		texture, texExists := rs.resourcesManager.GetTexture(entity.Sprite.TexturePath)

		if !texExists {

			// Attempt to load the texture if not already loaded
			_, err := rs.resourcesManager.LoadTexture(entity.Sprite.TexturePath)
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

	raylib.EndMode2D()
}

func (rs *RenderSystem) GetPriority() int {
	return rs.priority
}

func (rs *RenderSystem) SetCameraSystem(cs *CameraSystem) {
	rs.cameraSystem = cs
}
