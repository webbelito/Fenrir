package ecs

import (
	Raylib "github.com/gen2brain/raylib-go/raylib"
)

type Position struct {
	Vector Raylib.Vector2
}

type Velocity struct {
	Vector Raylib.Vector2
}

type Color struct {
	Color Raylib.Color
}

type Speed struct {
	Value float32
}
