package systems

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// ParticleSystem is a system that handles particles
type ParticleSystem struct {
	manager  *ecs.Manager
	priority int
}

// NewParticleSystem creates a new ParticleSystem
func NewParticleSystem(m *ecs.Manager, p int) *ParticleSystem {
	return &ParticleSystem{
		manager:  m,
		priority: p,
	}
}

func (ps *ParticleSystem) Update(dt float64) {

	// Get all entities with ParticleEmitterComponent
	entities := ps.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.ParticleEmitterComponent})

	// Iterate over all entities with ParticleEmitterComponent
	for _, entity := range entities {

		// Get the ParticleEmitterComponent
		emitterComp, emitterCompExists := ps.manager.GetComponent(entity, ecs.ParticleEmitterComponent)

		// Check if the ParticleEmitterComponent exists
		if !emitterCompExists {
			utils.ErrorLogger.Println("ParticleEmitterComponent does not exist")
			continue
		}

		// Cast the component to a ParticleEmitter
		emitter, ok := emitterComp.(*components.ParticleEmitter)

		// Check if the cast was successful
		if !ok {
			utils.ErrorLogger.Println("Failed to cast ParticleEmitterComponent to ParticleEmitter")
			continue
		}

		// Emit new particles
		if emitter.IsEmitting && time.Since(emitter.LastEmitTime) > time.Second/time.Duration(emitter.EmitRate) {

			// Reset the last emit time
			emitter.LastEmitTime = time.Now()

			// Create a new particle
			particle := &components.Particle{
				Position:     raylib.Vector2{X: 500, Y: 500},
				Velocity:     raylib.NewVector2(float32(raylib.GetRandomValue(-50, 50)), float32(raylib.GetRandomValue(-50, 50))),
				Acceleration: raylib.NewVector2(0, 0),
				Color:        raylib.Brown,
				Size:         5,
				Lifetime:     emitter.ParticleLifetime,
				Age:          0,
			}

			// Add the particle to the emitter
			emitter.Particles = append(emitter.Particles, particle)
		}

		// Create a new slice to store alive particles
		aliveParticles := []*components.Particle{}

		// Iterate over all particles
		for _, particle := range emitter.Particles {

			// Update particle age
			particle.Age += time.Duration(dt * float64(time.Second))

			// Check if the particle is alive
			if particle.Age < particle.Lifetime {
				// Update particle physics
				particle.Velocity = raylib.Vector2Add(particle.Velocity, raylib.Vector2Scale(particle.Acceleration, float32(dt)))
				particle.Position = raylib.Vector2Add(particle.Position, raylib.Vector2Scale(particle.Velocity, float32(dt)))
				aliveParticles = append(aliveParticles, particle)
			}
		}

		// Update the particles slice
		emitter.Particles = aliveParticles
	}
}

/*
GetPriority returns the priority of the system
*/
func (ps *ParticleSystem) GetPriority() int {
	return ps.priority
}
