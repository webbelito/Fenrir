package systeminterfaces

import "github.com/webbelito/Fenrir/pkg/components"

type Updatable interface {
	Update(dt float64)
	GetPriority() int
}

type Renderable interface {
	Render()
	GetPriority() int
}

type CameraSystemInterface interface {
	Update(dt float64)
	Render()
	GetCamera() *components.Camera
	SetCamera(camera *components.Camera)
	SetOwner(owner uint64)
	GetOwner() uint64
	GetPriority() int
}
