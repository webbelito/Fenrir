package systeminterfaces

import "github.com/webbelito/Fenrir/pkg/components"

// SystemInterface is the interface that all systems must implement
type UpdatableSystemInterface interface {
	Update(dt float64)
	GetPriority() int
}

// SystemInterface is the interface that all systems must implement
type RenderableSystemInterface interface {
	Render()
	GetPriority() int
}

// UIRenderableSystemInterface is the interface that all UI render systems must implement
type UIRenderableSystemInterface interface {
	RenderUI()
	GetPriority() int
}

// CameraSystemInterface is the interface that all camera systems must implement
type CameraSystemInterface interface {
	Update(dt float64)
	Render()
	GetCamera() *components.Camera
	SetCamera(camera *components.Camera)
	SetOwner(owner uint64)
	GetOwner() uint64
	GetPriority() int
}
