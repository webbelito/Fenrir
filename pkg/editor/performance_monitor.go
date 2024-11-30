package editor

import (
	"fmt"
	"time"

	"github.com/webbelito/Fenrir/pkg/ecs"
	metricinterface "github.com/webbelito/Fenrir/pkg/interfaces/metricinterfaces"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// historySize is the number of historical data points to store
const (
	historySize = 60
)

// PerformanceMonitor is a struct that represents the Performance Monitor
type PerformanceMonitor struct {
	manager                *ecs.Manager
	layoutPosition         raylib.Vector2
	performanceMonitorData PerformanceMonitorData
	historyIndex           int
}

// PerformanceMonitorData is a struct that holds the data for the Performance Monitor
type PerformanceMonitorData struct {
	FPS                int32
	performanceMetrics metricinterface.PerformanceMetrics

	// Historical data
	FPSHistory            [historySize]int32
	UpdateDurationHistory [historySize]time.Duration
	RenderDurationHistory [historySize]time.Duration
	TotalDurationHistory  [historySize]time.Duration
}

// NewPerformanceMonitor creates a new PerformanceMonitor
func NewPerformanceMonitor(m *ecs.Manager) *PerformanceMonitor {
	pm := &PerformanceMonitor{
		manager:                m,
		performanceMonitorData: PerformanceMonitorData{},
	}

	// Pre-fill FPSHistory with an inital FPS value e.g. 60
	for i := 0; i < historySize; i++ {
		pm.performanceMonitorData.FPSHistory[i] = 60
		pm.performanceMonitorData.UpdateDurationHistory[i] = 0
		pm.performanceMonitorData.RenderDurationHistory[i] = 0
		pm.performanceMonitorData.TotalDurationHistory[i] = 0
	}

	return pm

}

// SetPosition sets the position of the Performance Monitor
func (pm *PerformanceMonitor) SetPosition(position raylib.Vector2) {
	pm.layoutPosition = position
}

func (pm *PerformanceMonitor) Update() {

	// Get the performance metrics from the ECS Manager
	pm.performanceMonitorData.performanceMetrics = pm.manager.GetPerformanceMetrics()

	FPS := raylib.GetFPS()

	// Update the Performance Monitor data
	pm.performanceMonitorData.FPS = FPS

	// Store the historical data
	pm.performanceMonitorData.FPSHistory[pm.historyIndex] = FPS
	pm.performanceMonitorData.UpdateDurationHistory[pm.historyIndex] = pm.performanceMonitorData.performanceMetrics.UpdateDuration
	pm.performanceMonitorData.RenderDurationHistory[pm.historyIndex] = pm.performanceMonitorData.performanceMetrics.RenderDuration
	pm.performanceMonitorData.TotalDurationHistory[pm.historyIndex] = pm.performanceMonitorData.performanceMetrics.TotalDuration

	// Update the history index
	pm.historyIndex = (pm.historyIndex + 1) % historySize
}

func (pm *PerformanceMonitor) Render() {

	// Define the Performance Monitor panel position and size relative to the editor
	performanceMonitorRect := raylib.Rectangle{
		X:      pm.layoutPosition.X,
		Y:      pm.layoutPosition.Y,
		Width:  300,
		Height: 600,
	}

	// Draw the Performance Monitor panel
	raygui.Panel(performanceMonitorRect, "Performance Monitor")

	// Define spacing
	xOffset := performanceMonitorRect.X + 10
	yOffset := performanceMonitorRect.Y + 40
	lineHeight := 20
	lineSpacing := 30

	// Display Averaged Metrics
	averageFPS := pm.performanceMonitorData.calculateAverageFPS()
	averageUpdate := pm.performanceMonitorData.calculateAverageUpdateDuration()
	averageRender := pm.performanceMonitorData.calculateAverageRenderDuration()
	averageTotal := pm.performanceMonitorData.calculateAverageTotalDuration()

	// Convert time.Duration to microseconds
	averageUpdateMicro := averageUpdate.Microseconds()
	averageRenderMicro := averageRender.Microseconds()
	averageTotalMicro := averageTotal.Microseconds()

	// FPS Label
	fpsDisplayText := fmt.Sprintf("Avg FPS: %.2f", averageFPS)
	raygui.Label(raylib.Rectangle{
		X:      xOffset,
		Y:      yOffset,
		Width:  200,
		Height: float32(lineHeight),
	}, fpsDisplayText)

	// Update Duration Label
	updateDurationDisplayText := fmt.Sprintf("Avg Update: %d us", averageUpdateMicro)
	raygui.Label(raylib.Rectangle{
		X:      xOffset,
		Y:      yOffset + float32(lineSpacing),
		Width:  200,
		Height: float32(lineHeight),
	}, updateDurationDisplayText)

	// Render Duration Label
	renderDurationDisplayText := fmt.Sprintf("Avg Render: %d us", averageRenderMicro)
	raygui.Label(raylib.Rectangle{
		X:      xOffset,
		Y:      yOffset + float32(lineSpacing*2),
		Width:  200,
		Height: float32(lineHeight),
	}, renderDurationDisplayText)

	// Total Duration Label
	totalDurationDisplayText := fmt.Sprintf("Avg Total: %d us", averageTotalMicro)
	raygui.Label(raylib.Rectangle{
		X:      xOffset,
		Y:      yOffset + float32(lineSpacing*3),
		Width:  200,
		Height: float32(lineHeight),
	}, totalDurationDisplayText)

	// Draw FPS Graph
	drawGraph("FPS", pm.performanceMonitorData.FPSHistory[:], int32(performanceMonitorRect.X+10), int32(performanceMonitorRect.Y+180), 280, 100, raylib.Green)

	// Draw Update Duration Graph
	drawGraph("Update (us)", convertDurationsToInt32(pm.performanceMonitorData.UpdateDurationHistory[:]), int32(performanceMonitorRect.X+10), int32(performanceMonitorRect.Y+300), 280, 100, raylib.Red)
}

// calculateAverageFPS calculates the average FPS from the historical data
func (pmd *PerformanceMonitorData) calculateAverageFPS() float64 {

	// Calculate the sum of all FPS values
	var sum int32

	// Iterate over the FPS history
	for _, fps := range pmd.FPSHistory {
		sum += fps
	}

	// Calculate the average FPS
	return (float64(sum) / float64(historySize))
}

// calculateAverageUpdateDuration calculates the average Update Duration from the historical data
func (pmd *PerformanceMonitorData) calculateAverageUpdateDuration() time.Duration {

	// Calculate the sum of all Update Duration values
	var sum time.Duration

	// Iterate over the Update Duration history
	for _, updateDuration := range pmd.UpdateDurationHistory {
		sum += updateDuration
	}

	// Calculate the average Update Duration
	return sum / historySize
}

// calculateAverageRenderDuration calculates the average Render Duration from the historical data
func (pmd *PerformanceMonitorData) calculateAverageRenderDuration() time.Duration {

	// Calculate the sum of all Render Duration values
	var sum time.Duration

	// Iterate over the Render Duration history
	for _, renderDuration := range pmd.RenderDurationHistory {
		sum += renderDuration
	}

	// Calculate the average Render Duration
	return sum / historySize
}

// calculateAverageTotalDuration calculates the average Total Duration from the historical data
func (pmd *PerformanceMonitorData) calculateAverageTotalDuration() time.Duration {

	// Calculate the sum of all Total Duration values
	var sum time.Duration

	// Iterate over the Total Duration history
	for _, totalDuration := range pmd.TotalDurationHistory {
		sum += totalDuration
	}

	// Calculate the average Total Duration
	return sum / historySize
}

// Helper function to draw a simple line graph
func drawGraph(title string, data []int32, x int32, y int32, w int32, h int32, c raylib.Color) {

	// Draw title
	raygui.Label(raylib.Rectangle{
		X:      float32(x),
		Y:      float32(y - 20),
		Width:  float32(w),
		Height: 20,
	}, title)

	// Define graph area
	graphArea := raylib.Rectangle{
		X:      float32(x),
		Y:      float32(y),
		Width:  float32(w),
		Height: float32(h),
	}

	// Draw graph border
	raylib.DrawRectangleLinesEx(graphArea, 2, raylib.Gray)

	// Determine maxmimu, value dynamically
	var maxValue int32 = 100 // Default max value
	if title == "Update (us)" {
		maxValue = 2000 // Example max value for update duration
	}

	// Calculate the actual max from data
	for _, value := range data {
		if value > maxValue {
			maxValue = value
		}
	}

	// Prevent division by zero
	if maxValue == 0 {
		maxValue = 1
	}

	// Calculate vertical scale based on max value
	scaleY := float32(h) / float32(maxValue)

	// Calculate horizontal scale based on history size
	historySize := len(data)
	if historySize == 0 {
		historySize = 1
	}

	// Calculate the horizontal scale
	scaleX := float32(w) / float32(historySize)

	// Plot data points
	prevX := float32(x)
	prevY := float32(y+h) - (float32(data[0]) * scaleY)

	// Iterate over data points
	for i := 1; i < len(data); i++ {
		currentX := float32(x) + float32(i)*scaleX
		currentY := float32(y+h) - (float32(data[i]) * scaleY)

		// Draw the line
		raylib.DrawLineEx(raylib.Vector2{X: prevX, Y: prevY}, raylib.Vector2{X: currentX, Y: currentY}, 2, c)
		prevX = currentX
		prevY = currentY

	}
}

// Helper function to convert a slice of time.Duration to int64 (microseconds)
func convertDurationsToInt32(durations []time.Duration) []int32 {

	// Convert time.Duration to int32 (microseconds)
	result := make([]int32, len(durations))

	// Iterate over durations
	for i, d := range durations {
		result[i] = int32(d.Microseconds())
	}

	return result
}
