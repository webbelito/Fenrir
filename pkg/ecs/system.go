package ecs

import rl "github.com/gen2brain/raylib-go/raylib"

type System interface {
	Update(dt float64, em *EntitiesManager, cm *ComponentsManager)
}

type MovementSystem struct{}

func (ms *MovementSystem) Update(dt float64, em *EntitiesManager, cm *ComponentsManager) {

	positions, posExist := cm.Components[PositionComponent]
	velocities, velExist := cm.Components[VelocityComponent]

	if !posExist || !velExist {
		return
	}

	for entity, vel := range velocities {
		position, posExists := positions[entity].(*Position)
		velocity, velExists := vel.(*Velocity)

		if !posExists || !velExists {
			continue
		}

		deltaVelocity := rl.Vector2Scale(velocity.Vector, float32(dt))
		position.Vector = rl.Vector2Add(position.Vector, deltaVelocity)
	}
}
