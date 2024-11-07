// internal/ecs/systems/render_system.go

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/internal/ecs"
	"github.com/webbelito/Fenrir/internal/ecs/components"
)

type RenderSystem struct {
	windowWidth  int32
	windowHeight int32
}

func NewRenderSystem(windowWidth, windowHeight int32) *RenderSystem {
	return &RenderSystem{
		windowWidth:  windowWidth,
		windowHeight: windowHeight,
	}
}

func (rs *RenderSystem) Update(deltaTime float32, manager *ecs.Manager) {

	// Get all entities with a position and player component
	entities := manager.GetEntitiesWithComponents(&components.Position{}, &components.Player{})
	for _, entity := range entities {
		position := manager.GetComponent(entity, &components.Position{}).(*components.Position)
		player := manager.GetComponent(entity, &components.Player{}).(*components.Player)

		// Draw a circle representing the player
		rl.DrawCircle(int32(position.X), int32(position.Y), 10, rl.Red)

		// Draw the player's name above the circle
		rl.DrawText(player.Name, int32(position.X)-10, int32(position.Y)-20, 10, rl.Black)
	}

}
