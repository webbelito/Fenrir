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
	ID     uint64
	Vector raylib.Vector2
	Color  raylib.Color
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

	// Get all entities with a position component
	allPositionsComp, allPosExists := rs.componentsManager.Components[ecs.PositionComponent]

	// Get all entities with a color component
	allColorsComp, allColExists := rs.componentsManager.Components[ecs.ColorComponent]

	if !allPosExists && !allColExists {
		return
	}

	// Clear the entities slice while retaining capacity
	rs.Entities = rs.Entities[:0]

	// Collect entites that have both Position and Color components
	for entity, pos := range allPositionsComp {
		posComp, posCompExists := pos.(*components.Position)
		colorComp, colorCompExists := allColorsComp[entity].(*components.Color)

		if !posCompExists && !colorCompExists {
			continue
		}

		// Check if the entity is within the screen bounds
		if !raylib.CheckCollisionPointRec(posComp.Vector, rs.ScreenCullingRect) {
			continue
		}

		// Add the entity to the slice
		rs.Entities = append(rs.Entities, EntityData{
			ID:     entity.ID,
			Vector: posComp.Vector,
			Color:  colorComp.Color,
		})

	}

	// Sort entites by ID to ensure consistent rendering order
	/*sort.SliceStable(rs.Entities, func(i, j int) bool {
		return rs.Entities[i].ID < rs.Entities[j].ID
	})
	*/

	// Sort entities by ID to ensure consistent rendering order
	sort.SliceStable(rs.Entities, func(i int, j int) bool {
		return rs.Entities[i].ID < rs.Entities[j].ID // Not sure if ID.ID is the best way to do this
	})

	// Render entities
	for _, entity := range rs.Entities {
		raylib.DrawRectangleV(entity.Vector, raylib.NewVector2(5, 5), entity.Color)
	}
}
