package physicscomponents

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type RigidBody struct {
	Mass         float32
	InvMass      float32
	Velocity     raylib.Vector2
	Acceleration raylib.Vector2
	Force        raylib.Vector2
	Drag         float32
	Restitution  float32
	IsKinematic  bool
	IsStatic     bool
}

func NewRigidBody(mass float32, drag float32, restitution float32, isKinematic bool, isStatic bool) *RigidBody {
	rb := &RigidBody{
		Mass:         mass,
		Velocity:     raylib.NewVector2(0, 0),
		Acceleration: raylib.NewVector2(0, 0),
		Force:        raylib.NewVector2(0, 0),
		Drag:         drag,
		Restitution:  restitution,
		IsKinematic:  isKinematic,
		IsStatic:     isStatic,
	}

	// Calculate InvMass
	if isStatic || mass == 0 {
		rb.InvMass = 0
	} else {
		rb.InvMass = 1 / mass
	}

	return rb
}
