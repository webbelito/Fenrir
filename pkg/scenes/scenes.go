package scenes

import "github.com/webbelito/Fenrir/pkg/ecs"

type Scene interface {
	Initialize()
	Update(dt float64)
	Render()
	Cleanup()
	Pause()
	Resume()

	AddEntity(eID *ecs.Entity)
	RemoveAllEntities()
}
