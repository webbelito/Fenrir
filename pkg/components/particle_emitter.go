package components

import (
	"time"
)

type ParticleEmitter struct {
	Particles        []*Particle
	EmitRate         int
	ParticleLifetime time.Duration
	LastEmitTime     time.Time
	IsEmitting       bool
}
