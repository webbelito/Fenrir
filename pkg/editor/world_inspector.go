package editor

import (
	"fmt"

	"github.com/webbelito/Fenrir/pkg/ecs"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type WorldInspector struct {
	windowPosition raylib.Vector2
	ecsManager     *ecs.ECSManager
}

func NewWorldInspector(ecsM *ecs.ECSManager) *WorldInspector {
	return &WorldInspector{
		windowPosition: raylib.NewVector2(0, 0),
		ecsManager:     ecsM,
	}
}

func (wi *WorldInspector) SetPosition(position raylib.Vector2) {
	wi.windowPosition = position
}

func (wi *WorldInspector) Update() {
	// Currently no logic is needed for the World Inspector
}

func (wi *WorldInspector) Render() {

	// Define the World Inspector panel position and size relative to the editor
	inspectorRect := raylib.Rectangle{
		X:      wi.windowPosition.X,
		Y:      wi.windowPosition.Y,
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

	entityCount := wi.ecsManager.GetEntityCount()

	displayText := fmt.Sprintf("Active Entities: %d", entityCount)

	// Draw the entity count label
	raygui.Label(entityCountLabelRect, displayText)

}
