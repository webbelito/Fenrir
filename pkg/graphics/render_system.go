package graphics

import (
	"github.com/webbelito/Fenrir/pkg/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Render System is responsible for rendering entities
type RenderSystem struct{}

func (rs *RenderSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	positions, posExist := cm.Components[ecs.PositionComponent]

	if !posExist {
		return
	}

	for _, pos := range positions {
		position, posExists := pos.(*ecs.Position)

		if !posExists {
			continue
		}

		rl.DrawCircleV(position.Vector, 10, rl.Red)
	}
}
