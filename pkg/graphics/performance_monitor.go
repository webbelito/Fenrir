package graphics

import (
	"fmt"
	"strconv"
	"time"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type PerformanceMonitor struct {
	LayoutPosition         raylib.Vector2
	PerformanceMonitorData PerformanceMonitorData
}

type PerformanceMonitorData struct {
	FPS            int32
	UpdateDuration time.Duration
	RenderDuration time.Duration
	TotalDuration  time.Duration
}

func NewPerformanceMonitor(layoutPosition raylib.Vector2) *PerformanceMonitor {
	return &PerformanceMonitor{
		LayoutPosition:         layoutPosition,
		PerformanceMonitorData: PerformanceMonitorData{},
	}
}

func (pm *PerformanceMonitor) Update() {
	// Currently no logic is needed for the Performance Monitor
}

func (pm *PerformanceMonitor) Draw(pmd *PerformanceMonitorData) {

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

	// Display the FPS
	fpsDisplayText := "FPS: " + strconv.Itoa(int(pmd.FPS))

	// Draw the FPS label
	raygui.Label(fpsLabelRect, fpsDisplayText)

	// Define the position for the Update Duration label within the panel
	updateDurationLabelRect := raylib.Rectangle{
		X:      performanceMonitorRect.X + 10,
		Y:      performanceMonitorRect.Y + 70,
		Width:  200,
		Height: 20,
	}

	// Display the Update Duration
	updateDurationDisplayText := fmt.Sprintf("Update Duration: %d", pmd.UpdateDuration.Milliseconds())

	// Draw the Update Duration label
	raygui.Label(updateDurationLabelRect, updateDurationDisplayText)

	// Define the position for the Draw Duration label within the panel
	renderDurationLabelRect := raylib.Rectangle{
		X:      performanceMonitorRect.X + 10,
		Y:      performanceMonitorRect.Y + 100,
		Width:  200,
		Height: 20,
	}

	// Display the Render Duration
	renderDurationDisplayText := fmt.Sprintf("Render Duration: %d", pmd.RenderDuration.Milliseconds())

	// Draw the Draw Duration label
	raygui.Label(renderDurationLabelRect, renderDurationDisplayText)

	// Define the position for the Total Duration label within the panel
	totalDurationLabelRect := raylib.Rectangle{
		X:      performanceMonitorRect.X + 10,
		Y:      performanceMonitorRect.Y + 130,
		Width:  200,
		Height: 20,
	}

	// Display the Total Duration
	totalDurationDisplayText := fmt.Sprintf("Total Duration: %d", pmd.TotalDuration.Milliseconds())

	// Draw the Total Duration label
	raygui.Label(totalDurationLabelRect, totalDurationDisplayText)

}
