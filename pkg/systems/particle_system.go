package systems

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type ParticleSystem struct {
	ecsManager *ecs.ECSManager
	priority   int
}

func NewParticleSystem(ecsM *ecs.ECSManager, p int) *ParticleSystem {
	return &ParticleSystem{
		ecsManager: ecsM,
		priority:   p,
	}
}

func (ps *ParticleSystem) Update(dt float64) {
	entities := ps.ecsManager.GetComponentsManager().GetEntitiesWithComponents([]ecs.ComponentType{ecs.ParticleEmitterComponent})
	for _, entity := range entities {
		emitterComp, emitterCompExists := ps.ecsManager.GetComponent(entity, ecs.ParticleEmitterComponent)

		if !emitterCompExists {
			continue
		}

		emitter := emitterComp.(*components.ParticleEmitter)

		// Emit new particles
		if emitter.IsEmitting && time.Since(emitter.LastEmitTime) > time.Second/time.Duration(emitter.EmitRate) {
			emitter.LastEmitTime = time.Now()

			particle := &components.Particle{
				Position:     raylib.Vector2{X: 500, Y: 500},
				Velocity:     raylib.NewVector2(float32(raylib.GetRandomValue(-50, 50)), float32(raylib.GetRandomValue(-50, 50))),
				Acceleration: raylib.NewVector2(0, 0),
				Color:        raylib.Brown,
				Size:         5,
				Lifetime:     emitter.ParticleLifetime,
				Age:          0,
			}
			emitter.Particles = append(emitter.Particles, particle)
		}

		// Update particles
		aliveParticles := []*components.Particle{}
		for _, particle := range emitter.Particles {
			particle.Age += time.Duration(dt * float64(time.Second))

			if particle.Age < particle.Lifetime {
				// Update particle physics
				particle.Velocity = raylib.Vector2Add(particle.Velocity, raylib.Vector2Scale(particle.Acceleration, float32(dt)))
				particle.Position = raylib.Vector2Add(particle.Position, raylib.Vector2Scale(particle.Velocity, float32(dt)))
				aliveParticles = append(aliveParticles, particle)
			}
		}
		emitter.Particles = aliveParticles
	}
}

func (ps *ParticleSystem) GetPriority() int {
	return ps.priority
}
