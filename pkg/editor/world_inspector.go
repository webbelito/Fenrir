package editor

import (
	"fmt"

	"github.com/webbelito/Fenrir/pkg/ecs"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// WorldInspector is a struct that represents the World Inspector
type WorldInspector struct {
	windowPosition raylib.Vector2
	manager        *ecs.Manager
}

// NewWorldInspector creates a new WorldInspector
func NewWorldInspector(m *ecs.Manager) *WorldInspector {
	return &WorldInspector{
		windowPosition: raylib.NewVector2(0, 0),
		manager:        m,
	}
}

// SetPosition sets the position of the World Inspector
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

	// Get the entity count from the ECS Manager
	entityCount := wi.manager.GetEntityCount()

	// Create the display text for the entity count label
	displayText := fmt.Sprintf("Active Entities: %d", entityCount)

	// Draw the entity count label
	raygui.Label(entityCountLabelRect, displayText)

}
