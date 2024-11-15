package physicssystems

import (
	"math"

	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	physicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"
	"github.com/webbelito/Fenrir/pkg/utils"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type CollisionSystem struct {
	entitiesManager  *ecs.EntitiesManager
	componentManager *ecs.ComponentsManager
}

func NewCollisionSystem() *CollisionSystem {
	return &CollisionSystem{}
}

func (cs *CollisionSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	if em == nil || cm == nil {
		utils.ErrorLogger.Println("CollisionSystem: EntitiesManager or ComponentsManager is nil")
		return
	}

	// Assign the entities and components manager to the system
	cs.entitiesManager = em
	cs.componentManager = cm

	// Get all transform components
	transformComp, transformExists := cs.componentManager.Components[ecs.Transform2DComponent]

	// Get all rigid body components
	rigidBodyComps, rigidBodyCompsExist := cs.componentManager.Components[ecs.RigidBodyComponent]

	// Get all collider components
	colliderComps, colliderCompsExist := cs.componentManager.Components[ecs.BoxColliderComponent]

	if !transformExists || !rigidBodyCompsExist || !colliderCompsExist {
		return
	}

	// Get all entities
	entities := cs.componentManager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.Transform2DComponent, ecs.RigidBodyComponent, ecs.BoxColliderComponent})

	// Iterate over all entities unique pairs for collision detection
	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {

			// Get the entities
			entityA := entities[i]
			entityB := entities[j]

			// Get the position components
			posA := transformComp[entityA].(*components.Transform2D).Position
			posB := transformComp[entityB].(*components.Transform2D).Position

			// Get the rigid body components
			rbA := rigidBodyComps[entityA].(*physicscomponents.RigidBody)
			rbB := rigidBodyComps[entityB].(*physicscomponents.RigidBody)

			// Get the collider components
			// TODO: Implement a more generic Collider Type
			colliderA := colliderComps[entityA].(*physicscomponents.BoxCollider)
			colliderB := colliderComps[entityB].(*physicscomponents.BoxCollider)

			// Handle Box to Box collision
			if colliderA.Type == "Square" && colliderB.Type == "Square" {
				cs.handleBoxToBoxCollision(entityA, posA, rbA, colliderA, entityB, posB, rbB, colliderB)
			}

			// TODO: Handle Box to Circle collision
			// TODO: Handle Circle to Circle collision
		}
	}
}

func (cs *CollisionSystem) handleBoxToBoxCollision(eA uint64, pA raylib.Vector2, rbA *physicscomponents.RigidBody, cA *physicscomponents.BoxCollider, eB uint64, pB raylib.Vector2, rbB *physicscomponents.RigidBody, cB *physicscomponents.BoxCollider) {

	// Get the width and height of the entities
	widthA := cA.Size.X
	heightA := cA.Size.Y

	widthB := cB.Size.X
	heightB := cB.Size.Y

	// Calculate the distance between the two entities
	deltaX := pB.X - pA.X
	deltaY := pB.Y - pA.Y

	// Calculate the overlap on both axes
	overlapX := (widthA/2 + widthB/2) - float32(math.Abs(float64(deltaX)))
	overlapY := (heightA/2 + heightB/2) - float32(math.Abs(float64(deltaY)))

	// if there is an overlap, a collision has occurred
	if overlapX > 0 && overlapY > 0 {

		// Determine the axis of least penetration
		if overlapX < overlapY {

			// Collision on the X axis
			normal := raylib.NewVector2(0, 0)

			// Check if X is positive
			if deltaX > 0 {
				normal.X = 1
			} else {
				normal.X = -1
			}

			// Resolve overlap
			separation := raylib.Vector2Scale(normal, overlapX/2)
			pA = raylib.Vector2Subtract(pA, separation)
			pB = raylib.Vector2Add(pB, separation)

			// Update positions
			cs.componentManager.Components[ecs.Transform2DComponent][eA].(*components.Transform2D).Position = pA
			cs.componentManager.Components[ecs.Transform2DComponent][eB].(*components.Transform2D).Position = pB

			// Calculate relative velocity
			relativeVelocity := raylib.Vector2Subtract(rbB.Velocity, rbA.Velocity)
			velocityAlongNormal := raylib.Vector2DotProduct(relativeVelocity, normal)

			// Do not resolve if velocities are separating
			if velocityAlongNormal > 0 {
				return
			}

			// Calculate restitution
			restitution := float32(math.Min(float64(rbA.Restitution), float64(rbB.Restitution)))

			// Calculate impulse scalar
			impulseScalar := -(1 + restitution) * velocityAlongNormal
			impulseScalar /= 1/rbA.Mass + 1/rbB.Mass

			// Calculate impulse
			impulse := raylib.Vector2Scale(normal, impulseScalar)

			// Apply impulse to the entities
			rbA.Velocity = raylib.Vector2Subtract(rbA.Velocity, raylib.Vector2Scale(impulse, 1/rbA.Mass))
			rbB.Velocity = raylib.Vector2Add(rbB.Velocity, raylib.Vector2Scale(impulse, 1/rbB.Mass))

		} else {

			// Collision on the Y axis
			normal := raylib.NewVector2(0, 0)

			// Check if Y is positive
			if deltaY > 0 {
				normal.Y = 1
			} else {
				normal.Y = -1
			}

			// Resolve overlap
			separation := raylib.Vector2Scale(normal, overlapY/2)
			pA = raylib.Vector2Subtract(pA, separation)
			pB = raylib.Vector2Add(pB, separation)

			// Update positions
			cs.componentManager.Components[ecs.Transform2DComponent][eA].(*components.Transform2D).Position = pA
			cs.componentManager.Components[ecs.Transform2DComponent][eB].(*components.Transform2D).Position = pB

			// Calculate relative velocity
			relativeVelocity := raylib.Vector2Subtract(rbB.Velocity, rbA.Velocity)
			velocityAlongNormal := raylib.Vector2DotProduct(relativeVelocity, normal)

			// Do not resolve if velocities are separating
			if velocityAlongNormal > 0 {
				return
			}

			// Calculate restitution
			restitution := float32(math.Min(float64(rbA.Restitution), float64(rbB.Restitution)))

			// Calculate impulse scalar
			impulseScalar := -(1 + restitution) * velocityAlongNormal
			impulseScalar /= 1/rbA.Mass + 1/rbB.Mass

			// Calculate impulse
			impulse := raylib.Vector2Scale(normal, impulseScalar)

			// Apply impulse to the entities
			rbA.Velocity = raylib.Vector2Subtract(rbA.Velocity, raylib.Vector2Scale(impulse, 1/rbA.Mass))
			rbB.Velocity = raylib.Vector2Add(rbB.Velocity, raylib.Vector2Scale(impulse, 1/rbB.Mass))
		}
	}
}
