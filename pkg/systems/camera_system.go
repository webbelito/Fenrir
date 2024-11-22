package systems

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type CameraSystem struct {
	ecsManager *ecs.ECSManager
	camera     *components.Camera
	priority   int
}

func NewCameraSystem(ecsM *ecs.ECSManager, p int) *CameraSystem {
	return &CameraSystem{
		ecsManager: ecsM,
		camera: &components.Camera{
			OwnerEntity: 0,
			Target:      raylib.Vector2{X: 0, Y: 0},
			Offset:      raylib.Vector2{X: float32(raylib.GetScreenWidth()) / 2, Y: float32(raylib.GetScreenHeight()) / 2},
			Zoom:        1.0,
		},
		priority: p,
	}
}

func (cs *CameraSystem) Update(dt float64) {
	if cs.camera.OwnerEntity == 0 {
		return
	}

	transformComp, transformCompExists := cs.ecsManager.GetComponent(cs.camera.OwnerEntity, ecs.Transform2DComponent)
	if !transformCompExists {
		return
	}

	transform := transformComp.(*components.Transform2D)

	// Update the camera's target to the owner's position
	cs.camera.Target = transform.Position
}

func (cs *CameraSystem) Render() {
	// Do nothing
}

func (cs *CameraSystem) GetOwner() uint64 {
	return cs.camera.OwnerEntity
}

func (cs *CameraSystem) GetCameraTarget() raylib.Vector2 {
	return cs.camera.Target
}

func (cs *CameraSystem) GetCameraOffset() raylib.Vector2 {
	return cs.camera.Offset
}

func (cs *CameraSystem) GetCameraZoom() float32 {
	return cs.camera.Zoom
}

func (cs *CameraSystem) GetCamera() *components.Camera {
	return cs.camera
}

func (cs *CameraSystem) SetCamera(c *components.Camera) {
	cs.camera = c
}

func (cs *CameraSystem) SetOwner(e uint64) {
	cs.camera.OwnerEntity = e
}

func (cs *CameraSystem) GetPriority() int {
	return cs.priority
}
