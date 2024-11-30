package systems

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type ParticleRenderSystem struct {
	manager  *ecs.Manager
	priority int
}

func NewParticleRenderSystem(m *ecs.Manager, p int) *ParticleRenderSystem {
	return &ParticleRenderSystem{
		manager:  m,
		priority: p,
	}
}

func (prs *ParticleRenderSystem) Render() {

	// Get all entities with ParticleEmitterComponent
	entities := prs.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.ParticleEmitterComponent})

	// Iterate over all entities with ParticleEmitterComponent
	for _, entity := range entities {

		// Get the ParticleEmitterComponent
		emitterComp, exist := prs.manager.GetComponent(entity, ecs.ParticleEmitterComponent)

		// Check if the ParticleEmitterComponent exists
		if !exist {
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

		// Iterate over all particles in the emitter
		for _, particle := range emitter.Particles {

			// Calculate the alpha value based on the particle's age
			alpha := float32(1.0 - particle.Age.Seconds()/particle.Lifetime.Seconds())

			// Fade the particle's color based on the alpha value
			color := raylib.Fade(particle.Color, alpha)

			// Draw the particle
			raylib.DrawCircleV(particle.Position, particle.Size, color)
		}
	}
}

/*
GetPriority returns the priority of the system
*/
func (prs *ParticleRenderSystem) GetPriority() int {
	return prs.priority
}
