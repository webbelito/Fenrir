package editor

import (
	"github.com/webbelito/Fenrir/pkg/ecs"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Editor struct {
	Visible            bool
	WorldInspector     *WorldInspector
	PerformanceMonitor *PerformanceMonitor
}

var (
	wiPosition = raylib.NewVector2(10, 10)
	pmPosition = raylib.NewVector2(10, 550)
)

func NewEditor(world *ecs.ECSManager) *Editor {

	raygui.SetStyle(raylib.FontDefault, raygui.TEXT_SIZE, 20)

	return &Editor{
		Visible:            false,
		WorldInspector:     NewWorldInstructor(wiPosition, world),
		PerformanceMonitor: NewPerformanceMonitor(pmPosition),
	}
}

func (e *Editor) ToggleVisibility() {
	e.Visible = !e.Visible
}

func (e *Editor) Update(pmd *PerformanceMonitorData) {
	if e.Visible {
		// Handle editor specific logic here
		e.WorldInspector.Update()

		// Update the Performance Monitor
		e.PerformanceMonitor.Update(pmd)
	}
}

func (e *Editor) Draw(pmd *PerformanceMonitorData) {
	if e.Visible {

		// Here you can add more UI elements like buttons, text fields etc.
		// For the World Inspector, we'll integrate it separately

		// Draw the World Inspector
		e.WorldInspector.Draw()

		// Draw the Performance Monitor
		e.PerformanceMonitor.Draw(pmd)

	}
}
