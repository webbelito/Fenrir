package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type UIPanel struct {
	Title     string
	Bounds    raylib.Rectangle
	IsVisible bool
}

type UILabel struct {
	Label     string
	Bounds    raylib.Rectangle
	IsVisible bool
}
