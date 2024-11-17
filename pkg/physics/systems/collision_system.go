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

// Constants for collision resolution
const (
	slop              = 0.01 // Allowable penetration for collision resolution
	velocityThreshold = 0.1  // Minimum velocity for collision resolution before clamping
	DampingFactor     = 0.98 // Damping factor for to reduce oscillation
	PercentCorrection = 0.8  // Percentage of penetration to correct
)

type CollisionSystem struct {
	quadTree             *physics.QuadTree
	csMutex              sync.RWMutex
	ShouldRenderQuadTree bool
}

func NewCollisionSystem(b physics.Rectangle, c int32, mD int32, cD int32) *CollisionSystem {
	return &CollisionSystem{
		quadTree:             physics.NewQuadTree(b, c, mD, cD),
		ShouldRenderQuadTree: false,
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
	cs.quadTree = physics.NewQuadTree(boundry, cs.quadTree.Capacity, cs.quadTree.MaxDepth, cs.quadTree.CurrentDepth)
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

	// TODO: Remove the handleInput to the Input Manager
	// Handle input to toggle the QuadTree rendering
	if raylib.IsKeyPressed(raylib.KeyQ) {
		cs.ToggleQuadTreeRender()
	}
}

func (cs *CollisionSystem) Render(em *ecs.EntitiesManager, cm *ecs.ComponentsManager) {
	if !cs.ShouldRenderQuadTree {
		return
	}

	cs.csMutex.RLock()
	defer cs.csMutex.RUnlock()

	// Traverse the QuadTree and render the nodes
	cs.drawQuadTreeNode(cs.quadTree, cs.quadTree.CurrentDepth, raylib.Green)
}

func (cs *CollisionSystem) drawQuadTreeNode(qt *physics.QuadTree, cD int32, bC raylib.Color) {

	if qt == nil {
		return
	}

	// Draw the boundry of the QuadTree node
	raylib.DrawRectangleLines(
		int32(qt.Boundry.Position.X),
		int32(qt.Boundry.Position.Y),
		int32(qt.Boundry.Width),
		int32(qt.Boundry.Height),
		bC,
	)

	// Recursively draw the children
	if qt.Divided {
		cs.drawQuadTreeNode(qt.NW, cD+1, bC)
		cs.drawQuadTreeNode(qt.NE, cD+1, bC)
		cs.drawQuadTreeNode(qt.SW, cD+1, bC)
		cs.drawQuadTreeNode(qt.SE, cD+1, bC)
	}
}

func (cs *CollisionSystem) ToggleQuadTreeRender() {
	cs.ShouldRenderQuadTree = !cs.ShouldRenderQuadTree
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

	// Determine the axis of least penetration
	var normal raylib.Vector2
	var penetration float32

	// Determine the smallest overlap
	if overlapX < overlapY {

		// Resolve collision on the X axis
		if deltaX > 0 {
			normal = raylib.NewVector2(1, 0)
		} else {
			normal = raylib.NewVector2(-1, 0)
		}
		penetration = overlapX
	} else {

		// Resolve collision on the Y axis
		if deltaY > 0 {
			normal = raylib.NewVector2(0, 1)
		} else {
			normal = raylib.NewVector2(0, -1)
		}
		penetration = overlapY
	}

	// Relative velocity
	relativeVelocity := raylib.Vector2Subtract(rbB.Velocity, rbA.Velocity)
	velocityAlongNormal := raylib.Vector2DotProduct(relativeVelocity, normal)

	// Do not resolve if velocities are separating
	if velocityAlongNormal > 0 {
		return
	}

	// Calculate restitution (bounciness)
	restitution := math.Min(float64(rbA.Restitution), float64(rbB.Restitution))

	// Calculate impulse scalar
	impulseScalar := -(1 + float32(restitution)) * velocityAlongNormal
	impulseScalar /= (rbA.InvMass + rbB.InvMass)

	// Calculate impulse vector
	impulse := raylib.Vector2Scale(normal, impulseScalar)

	// Apply impulse to the entities
	rbA.Velocity = raylib.Vector2Subtract(rbA.Velocity, raylib.Vector2Scale(impulse, rbA.InvMass))
	rbB.Velocity = raylib.Vector2Add(rbB.Velocity, raylib.Vector2Scale(impulse, rbB.InvMass))

	// Positional correction to prevent sinking and jitter
	correctionMagnitude := float32(float32(math.Max(float64(penetration-slop), float64(0.0)))/(rbA.InvMass+rbB.InvMass)) * PercentCorrection
	correction := raylib.Vector2Scale(normal, correctionMagnitude)

	// Apply positional correction based on inverse mass, only if not static
	if rbA.InvMass != 0 {
		tA.Position = raylib.Vector2Subtract(tA.Position, raylib.Vector2Scale(correction, rbA.InvMass))
	}
	if rbB.InvMass != 0 {
		tb.Position = raylib.Vector2Add(tb.Position, raylib.Vector2Scale(correction, rbB.InvMass))
	}

	// Apply daming to reduce residual velocity and	prevent	oscillation.
	rbA.Velocity = raylib.Vector2Scale(rbA.Velocity, DampingFactor)
	rbB.Velocity = raylib.Vector2Scale(rbB.Velocity, DampingFactor)

	// Clamp velocities below the threshold to prevent jitter
	if raylib.Vector2Length(rbA.Velocity) < velocityThreshold {
		rbA.Velocity = raylib.NewVector2(0, 0)
	}
	if raylib.Vector2Length(rbB.Velocity) < velocityThreshold {
		rbB.Velocity = raylib.NewVector2(0, 0)
	}
}

// TODO: Implement a Destroy Entity
