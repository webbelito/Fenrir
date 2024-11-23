package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

// AudioSource represents an audio source component
type AudioSource struct {
	FilePath   string
	Volume     float32
	Sound      raylib.Sound
	IsLooping  bool
	ShouldPlay bool
}

// Audiolistener represents an audio listener component
type AudioListener struct {
	Position raylib.Vector2
	Rotation float32
}
