package components

import (
	"time"
)

// ParticleEmitter is a component that holds a particle emitter
type ParticleEmitter struct {
	Particles        []*Particle
	EmitRate         int
	ParticleLifetime time.Duration
	LastEmitTime     time.Time
	IsEmitting       bool
}
