package graphics

import (
	"fmt"

	"github.com/webbelito/Fenrir/pkg/ecs"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type WorldInspector struct {
	LayoutPosition raylib.Vector2
	World          *ecs.ECSManager
}

func NewWorldInstructor(layoutPosition raylib.Vector2, world *ecs.ECSManager) *WorldInspector {
	return &WorldInspector{
		LayoutPosition: layoutPosition,
		World:          world,
	}
}

func (wi *WorldInspector) Update() {
	// Currently no logic is needed for the World Inspector
}

func (wi *WorldInspector) Draw() {

	// Define the World Inspector panel position and size relative to the editor
	inspectorRect := raylib.Rectangle{
		X:      wi.LayoutPosition.X,
		Y:      wi.LayoutPosition.Y,
		Width:  300,
		Height: 600,
	}

	// Draw the World Inspector panel
	raygui.Panel(inspectorRect, "World Inspector")

	// Define the position for the entity count label witihn the panel
	entityCountLabelRect := raylib.Rectangle{
		X:      inspectorRect.X + 10,
		Y:      inspectorRect.Y + 40,
		Width:  200,
		Height: 20,
	}

	entityCount := wi.World.GetEntityCount()

	displayText := fmt.Sprintf("Active Entities: %d", entityCount)

	// Draw the entity count label
	raygui.Label(entityCountLabelRect, displayText)

}
