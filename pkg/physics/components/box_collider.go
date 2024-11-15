package physicscomponents

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type BoxCollider struct {
	Type string
	Size raylib.Vector2
}
