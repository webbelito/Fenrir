package physicscomponents

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type RigidBody struct {
	Mass         float32
	Velocity     raylib.Vector2
	Acceleration raylib.Vector2
	Force        raylib.Vector2
	Drag         float32
	Restitution  float32
	IsKinematic  bool
	IsStatic     bool
}
