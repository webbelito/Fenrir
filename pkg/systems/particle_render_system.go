package systems

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type ParticleRenderSystem struct {
	ecsManager *ecs.ECSManager
	priority   int
}

func NewParticleRenderSystem(ecsM *ecs.ECSManager, p int) *ParticleRenderSystem {
	return &ParticleRenderSystem{
		ecsManager: ecsM,
		priority:   p,
	}
}

func (prs *ParticleRenderSystem) Render() {
	entities := prs.ecsManager.GetComponentsManager().GetEntitiesWithComponents([]ecs.ComponentType{ecs.ParticleEmitterComponent})
	for _, entity := range entities {
		emitterComp, emitterCompExists := prs.ecsManager.GetComponent(entity, ecs.ParticleEmitterComponent)

		if !emitterCompExists {
			continue
		}

		emitter := emitterComp.(*components.ParticleEmitter)

		for _, particle := range emitter.Particles {
			alpha := float32(1.0 - particle.Age.Seconds()/particle.Lifetime.Seconds())
			color := raylib.Fade(particle.Color, alpha)
			raylib.DrawCircleV(particle.Position, particle.Size, color)
		}
	}
}

func (prs *ParticleRenderSystem) GetPriority() int {
	return prs.priority
}
