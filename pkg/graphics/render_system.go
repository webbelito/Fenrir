package graphics

import (
	"sort"

	"github.com/webbelito/Fenrir/pkg/ecs"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// Render System is responsible for rendering entities
type RenderSystem struct {
	ScreenCullingRect raylib.Rectangle
	Entities          []EntityData
}

type EntityData struct {
	ID     ecs.Entity
	Vector raylib.Vector2
	Color  raylib.Color
}

func NewRenderSystem(screenBounds raylib.Rectangle) *RenderSystem {
	return &RenderSystem{
		ScreenCullingRect: screenBounds,
		Entities:          []EntityData{},
	}
}

func (rs *RenderSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {
	// Currently the render systems doesn't update anything except rendering
}

func (rs *RenderSystem) Render(em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	// Get all entities with a position component
	allPositionsComp, allPosExists := cm.Components[ecs.PositionComponent]

	// Get all entities with a color component
	allColorsComp, allColExists := cm.Components[ecs.ColorComponent]

	if !allPosExists && !allColExists {
		return
	}

	// Clear the entities slice while retaining capacity
	rs.Entities = rs.Entities[:0]

	// Collect entites that have both Position and Color components
	for entity, pos := range allPositionsComp {
		posComp, posCompExists := pos.(*ecs.Position)
		colorComp, colorCompExists := allColorsComp[entity].(*ecs.Color)

		if !posCompExists && !colorCompExists {
			continue
		}

		// Check if the entity is within the screen bounds
		if !raylib.CheckCollisionPointRec(posComp.Vector, rs.ScreenCullingRect) {
			continue
		}

		rs.Entities = append(rs.Entities, EntityData{
			ID:     entity,
			Vector: posComp.Vector,
			Color:  colorComp.Color,
		})

	}

	// Sort entites by ID to ensure consistent rendering order
	sort.SliceStable(rs.Entities, func(i, j int) bool {
		return rs.Entities[i].ID < rs.Entities[j].ID
	})

	// Render entities
	for _, entity := range rs.Entities {
		raylib.DrawRectangleV(entity.Vector, raylib.NewVector2(5, 5), entity.Color)
	}

}
