package editor

import (
	"fmt"
	"time"

	raygui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	historySize = 60
)

type PerformanceMonitor struct {
	LayoutPosition         raylib.Vector2
	PerformanceMonitorData PerformanceMonitorData
	historyIndex           int
}

type PerformanceMonitorData struct {
	FPS            int32
	UpdateDuration time.Duration
	RenderDuration time.Duration
	TotalDuration  time.Duration

	// Historical data
	FPSHistory            [historySize]int32
	UpdateDurationHistory [historySize]time.Duration
	RenderDurationHistory [historySize]time.Duration
	TotalDurationHistory  [historySize]time.Duration
}

func NewPerformanceMonitor(layoutPosition raylib.Vector2) *PerformanceMonitor {
	pm := &PerformanceMonitor{
		LayoutPosition:         layoutPosition,
		PerformanceMonitorData: PerformanceMonitorData{},
	}

	// Pre-fill FPSHistory with an inital FPS value e.g. 60
	for i := 0; i < historySize; i++ {
		pm.PerformanceMonitorData.FPSHistory[i] = 60
		pm.PerformanceMonitorData.UpdateDurationHistory[i] = 0
		pm.PerformanceMonitorData.RenderDurationHistory[i] = 0
		pm.PerformanceMonitorData.TotalDurationHistory[i] = 0
	}

	return pm

}

func (pm *PerformanceMonitor) Update(pmd *PerformanceMonitorData) {

	// Update the Performance Monitor data
	pm.PerformanceMonitorData.FPS = pmd.FPS
	pm.PerformanceMonitorData.UpdateDuration = pmd.UpdateDuration
	pm.PerformanceMonitorData.RenderDuration = pmd.RenderDuration
	pm.PerformanceMonitorData.TotalDuration = pmd.TotalDuration

	// Store the historical data
	pm.PerformanceMonitorData.FPSHistory[pm.historyIndex] = pmd.FPS
	pm.PerformanceMonitorData.UpdateDurationHistory[pm.historyIndex] = pmd.UpdateDuration
	pm.PerformanceMonitorData.RenderDurationHistory[pm.historyIndex] = pmd.RenderDuration
	pm.PerformanceMonitorData.TotalDurationHistory[pm.historyIndex] = pmd.TotalDuration

	// Update the history index
	pm.historyIndex = (pm.historyIndex + 1) % historySize
}

func (pm *PerformanceMonitor) Draw(pmd *PerformanceMonitorData) {

	// Define the Performance Monitor panel position and size relative to the editor
	performanceMonitorRect := raylib.Rectangle{
		X:      pm.LayoutPosition.X,
		Y:      pm.LayoutPosition.Y,
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
	averageFPS := pm.PerformanceMonitorData.calculateAverageFPS()
	averageUpdate := pm.PerformanceMonitorData.calculateAverageUpdateDuration()
	averageRender := pm.PerformanceMonitorData.calculateAverageRenderDuration()
	averageTotal := pm.PerformanceMonitorData.calculateAverageTotalDuration()

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
	drawGraph("FPS", pm.PerformanceMonitorData.FPSHistory[:], int32(performanceMonitorRect.X+10), int32(performanceMonitorRect.Y+180), 280, 100, raylib.Green)

	// Draw Update Duration Graph
	drawGraph("Update (us)", convertDurationsToInt32(pm.PerformanceMonitorData.UpdateDurationHistory[:]), int32(performanceMonitorRect.X+10), int32(performanceMonitorRect.Y+300), 280, 100, raylib.Red)
}

func (pmd *PerformanceMonitorData) calculateAverageFPS() float64 {
	var sum int32
	for _, fps := range pmd.FPSHistory {
		sum += fps
	}

	avg := float64(sum) / float64(historySize)

	return avg
}

func (pmd *PerformanceMonitorData) calculateAverageUpdateDuration() time.Duration {
	var sum time.Duration
	for _, updateDuration := range pmd.UpdateDurationHistory {
		sum += updateDuration
	}
	return sum / historySize
}

func (pmd *PerformanceMonitorData) calculateAverageRenderDuration() time.Duration {
	var sum time.Duration
	for _, renderDuration := range pmd.RenderDurationHistory {
		sum += renderDuration
	}
	return sum / historySize
}

func (pmd *PerformanceMonitorData) calculateAverageTotalDuration() time.Duration {
	var sum time.Duration
	for _, totalDuration := range pmd.TotalDurationHistory {
		sum += totalDuration
	}
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

	scaleY := float32(h) / float32(maxValue)

	// Calculate horizontal scale based on history size
	historySize := len(data)
	if historySize == 0 {
		historySize = 1
	}

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
	result := make([]int32, len(durations))
	for i, d := range durations {
		result[i] = int32(d.Microseconds())
	}
	return result
}
