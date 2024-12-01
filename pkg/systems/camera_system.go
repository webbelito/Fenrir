package systems

import (
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// CameraSystem is a system that handles the camera
type CameraSystem struct {
	manager  *ecs.Manager
	camera   *components.Camera
	priority int
}

// NewCameraSystem creates a new CameraSystem
func NewCameraSystem(m *ecs.Manager, p int) *CameraSystem {
	cs := &CameraSystem{
		manager:  m,
		priority: p,
	}

	return cs

}

func (cs *CameraSystem) Update(dt float64) {

	cameraComp, exist := cs.manager.GetCamera()

	if !exist {
		utils.ErrorLogger.Println("CameraSystem: Camera component does not exist")
		return
	}

	cs.camera = cameraComp

	if cs.camera == nil {
		utils.ErrorLogger.Println("CameraSystem: Camera component is nil")
		return
	}

	// Get the transform component of the camera's owner entity
	transformComp, exist := cs.manager.GetComponent(cs.camera.OwnerEntity, ecs.Transform2DComponent)
	if !exist {
		return
	}

	// Cast the component to a Transform2D
	transform, ok := transformComp.(*components.Transform2D)

	// Check if the cast was successful
	if !ok {
		utils.ErrorLogger.Println("CameraSystem: Failed to cast Transform2DComponent to Transform2D")
		return
	}

	// Update the camera's target to the owner's position
	cs.camera.Target = transform.Position
}

func (cs *CameraSystem) Render() {
	// Do nothing
}

// GetOwner returns the owner entity of the camera
func (cs *CameraSystem) GetOwner() uint64 {
	return cs.camera.OwnerEntity
}

// GetCameraTarget returns the camera's target
func (cs *CameraSystem) GetCameraTarget() raylib.Vector2 {
	return cs.camera.Target
}

// GetCameraOffset returns the camera's offset
func (cs *CameraSystem) GetCameraOffset() raylib.Vector2 {
	return cs.camera.Offset
}

// GetCameraZoom returns the camera's zoom
func (cs *CameraSystem) GetCameraZoom() float32 {
	return cs.camera.Zoom
}

// SetCameraTarget sets the camera's target
func (cs *CameraSystem) GetCamera() *components.Camera {
	return cs.camera
}

// SetCameraTarget sets the camera's target
func (cs *CameraSystem) SetCamera(c *components.Camera) {
	cs.camera = c
}

// SetCameraTarget sets the camera's target
func (cs *CameraSystem) SetOwner(e uint64) {
	cs.camera.OwnerEntity = e
}

// GetPriority returns the priority of the system
func (cs *CameraSystem) GetPriority() int {
	return cs.priority
}
