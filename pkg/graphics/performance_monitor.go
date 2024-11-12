package graphics

import (
	"strconv"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type PerformanceMonitor struct {
	LayoutPosition raylib.Vector2
}

func NewPerformanceMonitor(layoutPosition raylib.Vector2) *PerformanceMonitor {
	return &PerformanceMonitor{
		LayoutPosition: layoutPosition,
	}
}

func (pm *PerformanceMonitor) Update() {
	// Currently no logic is needed for the Performance Monitor
}

func (pm *PerformanceMonitor) Draw() {

	// Define the Performance Monitor panel position and size relative to the editor
	performanceMonitorRect := raylib.Rectangle{
		X:      pm.LayoutPosition.X,
		Y:      pm.LayoutPosition.Y,
		Width:  300,
		Height: 300,
	}

	// Draw the Performance Monitor panel
	raygui.Panel(performanceMonitorRect, "Performance Monitor")

	// Define the position for the FPS label witihn the panel
	fpsLabelRect := raylib.Rectangle{
		X:      performanceMonitorRect.X + 10,
		Y:      performanceMonitorRect.Y + 40,
		Width:  200,
		Height: 20,
	}

	// Get the current FPS
	fps := int(raylib.GetFPS())

	displayText := "FPS: " + strconv.Itoa(fps)

	// Draw the FPS label
	raygui.Label(fpsLabelRect, displayText)
}
