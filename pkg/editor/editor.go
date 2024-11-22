package editor

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/ecs"
)

type Editor struct {
	worldInspector     *WorldInspector
	performanceMonitor *PerformanceMonitor
}

func NewEditor(ecsM *ecs.ECSManager) *Editor {
	ed := &Editor{
		worldInspector:     NewWorldInspector(ecsM),
		performanceMonitor: NewPerformanceMonitor(ecsM),
	}

	ed.setupLayouts()

	return ed
}

func (e *Editor) setupLayouts() {

	// Define positions within the editor window
	wiPosition := raylib.NewVector2(10, 10)
	pmPosition := raylib.NewVector2(10, 550)

	// Set the positions for each view
	e.worldInspector.SetPosition(wiPosition)
	e.performanceMonitor.SetPosition(pmPosition)
}

func (e *Editor) Update(dt float64) {
	e.worldInspector.Update()
	e.performanceMonitor.Update()
}

func (e *Editor) Render() {
	e.worldInspector.Render()
	e.performanceMonitor.Render()
}
