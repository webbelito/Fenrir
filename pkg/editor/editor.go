package editor

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/graphics"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Editor struct {
	Visible            bool
	WorldInspector     *graphics.WorldInspector
	PerformanceMonitor *graphics.PerformanceMonitor
}

var (
	wiPosition = raylib.NewVector2(10, 10)
	pmPosition = raylib.NewVector2(10, 550)
)

func NewEditor(world *ecs.ECSManager) *Editor {

	raygui.SetStyle(raylib.FontDefault, raygui.TEXT_SIZE, 20)

	return &Editor{
		Visible:            false,
		WorldInspector:     graphics.NewWorldInstructor(wiPosition, world),
		PerformanceMonitor: graphics.NewPerformanceMonitor(pmPosition),
	}
}

func (e *Editor) ToggleVisibility() {
	e.Visible = !e.Visible
}

func (e *Editor) Update() {
	if e.Visible {
		// Handle editor specific logic here
		e.WorldInspector.Update()
	}
}

func (e *Editor) Draw(pmd *graphics.PerformanceMonitorData) {
	if e.Visible {

		// Draw a semi-transparent background
		/*raygui.Panel(raylib.Rectangle{
			X:      10,
			Y:      10,
			Width:  float32(raylib.GetScreenWidth()) - 20,
			Height: float32(raylib.GetScreenHeight()) - 20,
		},
			"Fenrir Editor")
		*/

		// Here you can add more UI elements like buttons, text fields etc.
		// For the World Inspector, we'll integrate it separately

		// Draw the World Inspector
		e.WorldInspector.Draw()

		// Draw the Performance Monitor
		e.PerformanceMonitor.Draw(pmd)

	}
}
