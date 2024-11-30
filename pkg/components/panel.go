package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// UIPanel is a component that holds a UI panel
type UIPanel struct {
	Title     string
	Bounds    raylib.Rectangle
	IsVisible bool
}
