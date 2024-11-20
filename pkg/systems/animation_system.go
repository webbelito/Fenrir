package systems

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
)

type AnimationSystem struct {
	ecsManager *ecs.ECSManager
	priority   int
}

func NewAnimationSystem(ecsM *ecs.ECSManager, p int) *AnimationSystem {
	return &AnimationSystem{
		ecsManager: ecsM,
		priority:   p,
	}
}

func (as *AnimationSystem) Update(dt float64) {
	entities := as.ecsManager.GetComponentsManager().GetEntitiesWithComponents([]ecs.ComponentType{ecs.AnimationComponent})

	for _, entity := range entities {
		animationComp, animationCompExists := as.ecsManager.GetComponent(entity, ecs.AnimationComponent)

		if !animationCompExists {
			continue
		}

		animation := animationComp.(*components.Animation)

		if !animation.IsPlaying {
			continue
		}

		animation.ElapsedTime += time.Duration(dt * float64(time.Second))
		if animation.ElapsedTime >= animation.FrameDuration {
			animation.ElapsedTime = 0
			animation.CurrentFrame++
			if animation.CurrentFrame >= len(animation.Frames) {
				if animation.IsLooping {
					animation.CurrentFrame = 0
				} else {
					animation.IsPlaying = false
				}
			}
		}

		// Update the sprite's source rectangle
		spriteComp, spriteCompExists := as.ecsManager.GetComponent(entity, ecs.SpriteComponent)

		if !spriteCompExists {
			continue
		}

		sprite := spriteComp.(*components.Sprite)

		sprite.SourceRect = animation.Frames[animation.CurrentFrame]
	}
}

func (as *AnimationSystem) GetPriority() int {
	return as.priority
}
