package systems

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/resources"
)

type AudioSystem struct {
	ecsManager      *ecs.ECSManager
	resourceManager *resources.ResourcesManager
	priority        int
}

func NewAudioSystem(ecsM *ecs.ECSManager, rm *resources.ResourcesManager, p int) *AudioSystem {
	return &AudioSystem{
		ecsManager:      ecsM,
		resourceManager: rm,
		priority:        p,
	}
}

func (as *AudioSystem) Update(dt float64) {
	audioEntities := as.ecsManager.GetComponentsManager().GetEntitiesWithComponents([]ecs.ComponentType{ecs.AudioSourceComponent})
	for _, entity := range audioEntities {
		audioComp, audioCompExists := as.ecsManager.GetComponent(entity, ecs.AudioSourceComponent)

		if !audioCompExists {
			continue
		}

		audio := audioComp.(*components.AudioSource)

		if audio.ShouldPlay {
			raylib.PlaySound(audio.Sound)
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

func (as *AudioSystem) GetPriority() int {
	return as.priority
}
