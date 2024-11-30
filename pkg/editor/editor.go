package editor

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/ecs"
)

// Editor is a struct that represents the Editor
type Editor struct {
	worldInspector     *WorldInspector
	performanceMonitor *PerformanceMonitor
}

// NewEditor creates a new Editor
func NewEditor(m *ecs.Manager) *Editor {

	ed := &Editor{
		worldInspector:     NewWorldInspector(m),
		performanceMonitor: NewPerformanceMonitor(m),
	}

	ed.setupLayouts()

	return ed
}

func (e *Editor) Update(dt float64) {
	e.worldInspector.Update()
	e.performanceMonitor.Update()
}

func (e *Editor) Render() {
	e.worldInspector.Render()
	e.performanceMonitor.Render()
}

// setupLayouts sets up the layout for the editor
func (e *Editor) setupLayouts() {

	// Define positions within the editor window
	wiPosition := raylib.NewVector2(10, 10)
	pmPosition := raylib.NewVector2(10, 550)

	// Set the positions for each view
	e.worldInspector.SetPosition(wiPosition)
	e.performanceMonitor.SetPosition(pmPosition)
}
