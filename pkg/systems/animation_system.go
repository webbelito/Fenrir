package systems

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"
)

type AnimationSystem struct {
	manager  *ecs.Manager
	priority int
}

func NewAnimationSystem(m *ecs.Manager, p int) *AnimationSystem {
	return &AnimationSystem{
		manager:  m,
		priority: p,
	}
}

func (as *AnimationSystem) Update(dt float64) {

	// Get all entities with AnimationComponent
	entities := as.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.AnimationComponent})

	// Iterate over all entities with AnimationComponent
	for _, entity := range entities {

		// Get the AnimationComponent
		animationComp, exist := as.manager.GetComponent(entity, ecs.AnimationComponent)

		// Check if the AnimationComponent exists
		if !exist {
			utils.ErrorLogger.Println("AnimationComponent does not exist")
			continue
		}

		// Cast the component to an Animation
		animation, ok := animationComp.(*components.Animation)

		// Check if the cast was successful
		if !ok {
			utils.ErrorLogger.Println("Failed to cast AnimationComponent to Animation")
			continue
		}

		// Check if the animation is playing
		if !animation.IsPlaying {
			continue
		}

		// Update the animation
		animation.ElapsedTime += time.Duration(dt * float64(time.Second))

		// Check if the current frame duration has been reached
		if animation.ElapsedTime >= animation.FrameDuration {

			// Reset the elapsed time
			animation.ElapsedTime = 0
			animation.CurrentFrame++

			// Check if the animation has reached the end
			if animation.CurrentFrame >= len(animation.Frames) {

				// Check if the animation is looping
				if animation.IsLooping {
					animation.CurrentFrame = 0
				} else {
					animation.IsPlaying = false
				}
			}
		}

		// Update the sprite's source rectangle
		spriteComp, exist := as.manager.GetComponent(entity, ecs.SpriteComponent)

		if !exist {
			utils.ErrorLogger.Println("SpriteComponent does not exist")
			continue
		}

		// Cast the component to a Sprite
		sprite, ok := spriteComp.(*components.Sprite)

		if !ok {
			utils.ErrorLogger.Println("Failed to cast SpriteComponent to Sprite")
			continue
		}

		// Update the sprite's source rectangle
		sprite.SourceRect = animation.Frames[animation.CurrentFrame]
	}
}

// GetPriority returns the priority of the system
func (as *AnimationSystem) GetPriority() int {
	return as.priority
}
