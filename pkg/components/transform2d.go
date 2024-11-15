package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Transform2D struct {
	Position raylib.Vector2
	Rotation float32
	Scale    raylib.Vector2
}
