package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// UILabel is a component that holds a UI label
type UILabel struct {
	Label     string
	Bounds    raylib.Rectangle
	IsVisible bool
}
