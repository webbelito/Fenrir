package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// Camera is a component that holds a camera
type Camera struct {
	OwnerEntity uint64
	Target      raylib.Vector2
	Offset      raylib.Vector2
	Zoom        float32
}
