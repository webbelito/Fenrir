package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// Velocity is a component that holds a velocity vector
type Velocity struct {
	Vector raylib.Vector2
}
