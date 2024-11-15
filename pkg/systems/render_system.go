package systems

import (
	"sort"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type RenderSystem struct {
	ScreenCullingRect raylib.Rectangle
	Entities          []EntityData
	entitiesManager   *ecs.EntitiesManager
	componentsManager *ecs.ComponentsManager
}

type EntityData struct {
	ID       uint64
	Position raylib.Vector2
	Rotation float32
	Scale    raylib.Vector2
	Color    raylib.Color
}

func NewRenderSystem(screenBounds raylib.Rectangle) *RenderSystem {
	return &RenderSystem{
		ScreenCullingRect: screenBounds,
		Entities:          []EntityData{},
	}
}

func (rs *RenderSystem) Render(em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {
	if em == nil || cm == nil {
		utils.ErrorLogger.Println("RenderSystem: EntitiesManager or ComponentsManager is nil")
		return
	}

	// Assign the entities and components manager to the system
	rs.entitiesManager = em
	rs.componentsManager = cm

	rs.RenderEntities()

}

// TODO: Refactor this to not require an update method
func (rs *RenderSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {
	// Do nothing
}

func (rs *RenderSystem) RenderEntities() {

	transformComps, transformCompsExists := rs.componentsManager.Components[ecs.Transform2DComponent]
	colorComps, colorCompsExists := rs.componentsManager.Components[ecs.ColorComponent]

	if !transformCompsExists || !colorCompsExists {
		return
	}

	// Clear the entities slice while retaining capacity
	rs.Entities = rs.Entities[:0]

	// Collect entites that have both Position and Color components
	for entity, pos := range transformComps {
		transform, _ := pos.(*components.Transform2D)
		colorComp, colorCompExists := colorComps[entity].(*components.Color)

		if !colorCompExists {
			continue
		}

		// Check if the entity is within the screen bounds
		if !raylib.CheckCollisionPointRec(transform.Position, rs.ScreenCullingRect) {
			continue
		}

		// Add the entity to the slice
		rs.Entities = append(rs.Entities, EntityData{
			ID:       entity.ID,
			Position: transform.Position,
			Rotation: transform.Rotation,
			Scale:    transform.Scale,
			Color:    colorComp.Color,
		})

	}

	// Sort entities by ID to ensure consistent rendering order
	sort.SliceStable(rs.Entities, func(i int, j int) bool {
		return rs.Entities[i].ID < rs.Entities[j].ID
	})

	// Render entities
	for _, entity := range rs.Entities {
		raylib.DrawRectangleV(entity.Position, entity.Scale, entity.Color)
	}
}
