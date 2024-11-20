package components

import (
	"time"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Animation struct {
	Frames        []raylib.Rectangle
	CurrentFrame  int
	FrameDuration time.Duration
	ElapsedTime   time.Duration
	IsLooping     bool
	IsPlaying     bool
}
