package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Sprite struct {
	TexturePath string
	SourceRect  raylib.Rectangle
	DestRect    raylib.Rectangle
	Origin      raylib.Vector2
	Rotation    float32
	Color       raylib.Color
}
