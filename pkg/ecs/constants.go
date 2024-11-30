package ecs

// Component types for the ECS
const (
	Transform2DComponent = iota
	RigidBodyComponent
	BoxColliderComponent
	ColorComponent
	SpeedComponent
	PlayerComponent
	VelocityComponent
	SpriteComponent
	ParticleComponent
	ParticleEmitterComponent
	AnimationComponent
	AudioSourceComponent
	CameraComponent
)

// Component types for UI components with an offset of 100
const (
	UILabelComponent = iota + 100
	UIPanelComponent
	UIButtonComponent
)
