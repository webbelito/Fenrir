package components

import (
	"time"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// Particle is a component that holds a particle
type Particle struct {
	Position     raylib.Vector2
	Velocity     raylib.Vector2
	Acceleration raylib.Vector2
	Color        raylib.Color
	Size         float32
	Lifetime     time.Duration
	Age          time.Duration
}
