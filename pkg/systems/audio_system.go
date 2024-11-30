package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"
)

// AudioSystem is a system that handles audio
type AudioSystem struct {
	manager  *ecs.Manager
	priority int
}

// NewAudioSystem creates a new AudioSystem
func NewAudioSystem(m *ecs.Manager, p int) *AudioSystem {
	return &AudioSystem{
		manager:  m,
		priority: p,
	}
}

func (as *AudioSystem) Update(dt float64) {

	// Get all entities with AudioSourceComponent
	audioEntities := as.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.AudioSourceComponent})

	// Iterate over all entities with AudioSourceComponent
	for _, entity := range audioEntities {

		// Get the AudioSourceComponent
		audioComp, exist := as.manager.GetComponent(entity, ecs.AudioSourceComponent)

		// Check if the AudioSourceComponent exists
		if !exist {
			utils.ErrorLogger.Println("AudioSourceComponent does not exist")
			continue
		}

		// Cast the component to an AudioSource
		audio, ok := audioComp.(*components.AudioSource)

		// Check if the cast was successful
		if !ok {
			utils.ErrorLogger.Println("Failed to cast AudioSourceComponent to AudioSource")
			continue
		}

		// Play sound
		if audio.ShouldPlay {
			raylib.PlaySound(audio.Sound)

			// Reset ShouldPlay
			audio.ShouldPlay = false
		}

		// Adjust volume
		raylib.SetSoundVolume(audio.Sound, audio.Volume)

		// Handle looping
		if audio.IsLooping && !raylib.IsSoundPlaying(audio.Sound) {
			raylib.PlaySound(audio.Sound)
		}
	}
}

// GetPriority returns the priority of the system
func (as *AudioSystem) GetPriority() int {
	return as.priority
}
