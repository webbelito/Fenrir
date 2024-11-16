package physicssystems

import (
	// STD
	"math"
	"sync"

	// ECS
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"

	// UTILS
	"github.com/webbelito/Fenrir/pkg/utils"

	// PHYSICS
	"github.com/webbelito/Fenrir/pkg/physics"
	physicscomponents "github.com/webbelito/Fenrir/pkg/physics/components"

	// RAYLIB
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type CollisionSystem struct {
	quadTree *physics.QuadTree
	csMutex  sync.RWMutex
}

func NewCollisionSystem(b physics.Rectangle, c int32) *CollisionSystem {
	return &CollisionSystem{
		quadTree: physics.NewQuadTree(b, c),
	}
}

func (cs *CollisionSystem) Update(dt float64, em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {

	if em == nil || cm == nil {
		utils.ErrorLogger.Println("CollisionSystem: EntitiesManager or ComponentsManager is nil")
		return
	}

	// Define the boundry of the QuadTree based on screen size or game world
	// TODO: Implement a game world size for this

	screenWidth := float32(raylib.GetScreenWidth())
	screenHeight := float32(raylib.GetScreenHeight())

	boundry := physics.Rectangle{
		Position: raylib.NewVector2(0, 0),
		Width:    screenWidth,
		Height:   screenHeight,
	}

	// Initialize or reset the QuadTree
	cs.csMutex.Lock()
	cs.quadTree.Clear()
	cs.quadTree = physics.NewQuadTree(boundry, cs.quadTree.Capacity)
	cs.csMutex.Unlock()

	// Retrieve entities with the required components
	entities := cm.GetEntitiesWithComponents([]ecs.ComponentType{
		ecs.Transform2DComponent,
		ecs.BoxColliderComponent,
	})

	// Insert entities into the QuadTree
	for _, entity := range entities {
		transformComp, transCompExists := cm.GetComponent(entity, ecs.Transform2DComponent)
		if !transCompExists {
			continue
		}

		transform, transExists := transformComp.(*components.Transform2D)

		if !transExists {
			continue
		}

		// Insert entity into the QuadTree based on its position
		cs.csMutex.Lock()
		cs.quadTree.Insert(entity, transform.Position)
		cs.csMutex.Unlock()

	}

	// Perform collision checks
	for _, eID := range entities {

		// Get Transform2D component
		transformComp, transCompExists := cm.GetComponent(eID, ecs.Transform2DComponent)
		if !transCompExists {
			continue
		}

		// Get the reference to the Transform2D component
		transform, transExists := transformComp.(*components.Transform2D)
		if !transExists {
			continue
		}

		// Get the BoxCollider component
		colliderComp, colliderCompExists := cm.GetComponent(eID, ecs.BoxColliderComponent)
		if !colliderCompExists {
			continue
		}

		// Get the reference to the BoxCollider component
		collider, colliderExists := colliderComp.(*physicscomponents.BoxCollider)
		if !colliderExists {
			continue
		}

		// Define the range for querying the potential colliders
		rangeRect := physics.Rectangle{
			Position: raylib.NewVector2(
				transform.Position.X-collider.Size.X,
				transform.Position.Y-collider.Size.Y,
			),
			Width:  collider.Size.X * 2,
			Height: collider.Size.Y * 2,
		}

		// Query the QuadTree for potential colliders
		found := []uint64{}
		cs.csMutex.RLock()
		cs.quadTree.Query(rangeRect, &found)
		cs.csMutex.RUnlock()

		// Check for actiual collisions with the found entities
		for _, otherID := range found {
			if otherID == eID {
				continue // Ignore self
			}

			// Check for Box to Box collision
			cs.handleBoxToBoxCollision(eID, otherID, cm)
		}
	}
}

func (cs *CollisionSystem) handleBoxToBoxCollision(eA uint64, eB uint64, cm *ecs.ComponentsManager) {

	// Get the components for entity A
	transformA, transformAExists := cm.GetComponent(eA, ecs.Transform2DComponent)
	colliderA, colliderAExists := cm.GetComponent(eA, ecs.BoxColliderComponent)
	rigidBodyA, rigidbodyAExists := cm.GetComponent(eA, ecs.RigidBodyComponent)

	if !transformAExists || !colliderAExists || !rigidbodyAExists {
		return
	}

	// Get the components for entity B
	transformB, transformBExists := cm.GetComponent(eB, ecs.Transform2DComponent)
	colliderB, colliderBExists := cm.GetComponent(eB, ecs.BoxColliderComponent)
	rigidBodyB, rigidBodyBExists := cm.GetComponent(eB, ecs.RigidBodyComponent)

	if !transformBExists || !colliderBExists || !rigidBodyBExists {
		return
	}

	// Type assertions
	tA, tAOk := transformA.(*components.Transform2D)
	tb, tBOk := transformB.(*components.Transform2D)
	cA, cAOk := colliderA.(*physicscomponents.BoxCollider)
	cB, cBOk := colliderB.(*physicscomponents.BoxCollider)
	rbA, rbAOk := rigidBodyA.(*physicscomponents.RigidBody)
	rbB, rbBOk := rigidBodyB.(*physicscomponents.RigidBody)

	if !tAOk || !tBOk || !cAOk || !cBOk || !rbAOk || !rbBOk {
		return
	}

	// Calculate the difference in positions
	deltaX := tb.Position.X - tA.Position.X
	deltaY := tb.Position.Y - tA.Position.Y

	// Calculate the combined half-widths and half-heights
	halfWidthA := cA.Size.X / 2
	halfHeightA := cA.Size.Y / 2
	halfWidthB := cB.Size.X / 2
	halfHeightB := cB.Size.Y / 2

	// Calculate the overlap on both axes
	overlapX := (halfWidthA + halfWidthB) - float32(math.Abs(float64(deltaX)))
	overlapY := (halfHeightA + halfHeightB) - float32(math.Abs(float64(deltaY)))

	// If there's no overlap, no collision
	if overlapX <= 0 || overlapY <= 0 {
		return
	}

	// Determine the smallest overlap
	if overlapX < overlapY {

		// Resolve collision on the X axis
		if deltaX > 0 {
			tA.Position.X -= overlapX / 2
			tb.Position.X += overlapX / 2
		} else {
			tA.Position.X += overlapX / 2
			tb.Position.X -= overlapX / 2
		}

		// Adjust velocities based on restitution
		vxA := rbA.Velocity.X
		vxB := rbB.Velocity.X
		rX := (rbA.Restitution + rbB.Restitution) / 2

		// Swap velocities with restitution Forumal: v1' = (v1(m1-m2) + 2*m2*v2) / (m1 + m2) * r
		rbA.Velocity.X = (vxA*(rbA.Mass-rbB.Mass) + 2*rbB.Mass*vxB) / (rbA.Mass + rbB.Mass) * rX
		rbB.Velocity.X = (vxB*(rbB.Mass-rbA.Mass) + 2*rbA.Mass*vxA) / (rbA.Mass + rbB.Mass) * rX

	} else {

		// Resolve collision on the Y axis
		if deltaY > 0 {
			tA.Position.Y -= overlapY / 2
			tb.Position.Y += overlapY / 2
		} else {
			tA.Position.Y += overlapY / 2
			tb.Position.Y -= overlapY / 2
		}

		// Adjust velocities based on restitution
		vyA := rbA.Velocity.Y
		vyB := rbB.Velocity.Y
		rY := (rbA.Restitution + rbB.Restitution) / 2

		// Swap velocities with restitution Forumal: v1' = (v1(m1-m2) + 2*m2*v2) / (m1 + m2) * r
		rbA.Velocity.Y = (vyA*(rbA.Mass-rbB.Mass) + 2*rbB.Mass*vyB) / (rbA.Mass + rbB.Mass) * rY
		rbB.Velocity.Y = (vyB*(rbB.Mass-rbA.Mass) + 2*rbA.Mass*vyA) / (rbA.Mass + rbB.Mass) * rY

	}
}

// TODO: Implement a Destroy Entity

/*
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
*/
