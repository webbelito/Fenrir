package physics

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/utils"
)

// Recatangle defines the boundry of each QuadTree node
type Rectangle struct {
	Position raylib.Vector2
	Width    float32
	Height   float32
}

// Contains checks if a point is within the bounds of a rectangle
func (r *Rectangle) Contains(point raylib.Vector2) bool {
	return point.X >= r.Position.X &&
		point.X <= r.Position.X+r.Width &&
		point.Y >= r.Position.Y &&
		point.Y <= r.Position.Y+r.Height
}

// Intersects checks if two rectrangles intersect
func (r *Rectangle) Intersects(other *Rectangle) bool {
	return !(other.Position.X > r.Position.X+r.Width ||
		other.Position.X+other.Width < r.Position.X ||
		other.Position.Y > r.Position.Y+r.Height ||
		other.Position.Y+other.Height < r.Position.Y)
}

// QuadTree represents a node in the QuadTree
type QuadTree struct {
	Boundry  Rectangle
	Capacity int32
	Entities []uint64
	Divided  bool
	NE       *QuadTree
	NW       *QuadTree
	SE       *QuadTree
	SW       *QuadTree
}

// NewQuadTree initializes a new QuadTree node
func NewQuadTree(b Rectangle, c int32) *QuadTree {
	return &QuadTree{
		Boundry:  b,
		Capacity: c,
		Entities: make([]uint64, 0, c),
		Divided:  false,
	}
}

// Subdivide splits the QuadTree into four childen QuadTrees
func (qt *QuadTree) Subdivide() {
	halfWidth := qt.Boundry.Width / 2
	halfHeight := qt.Boundry.Height / 2
	x := qt.Boundry.Position.X
	y := qt.Boundry.Position.Y

	qt.NE = NewQuadTree(Rectangle{raylib.NewVector2(x+halfWidth, y), halfWidth, halfHeight}, qt.Capacity)
	qt.NW = NewQuadTree(Rectangle{raylib.NewVector2(x, y), halfWidth, halfHeight}, qt.Capacity)
	qt.SE = NewQuadTree(Rectangle{raylib.NewVector2(x+halfWidth, y+halfHeight), halfWidth, halfHeight}, qt.Capacity)
	qt.SW = NewQuadTree(Rectangle{raylib.NewVector2(x, y+halfHeight), halfWidth, halfHeight}, qt.Capacity)

	qt.Divided = true
}

// Insert adds an entity to the QuadTree
func (qt *QuadTree) Insert(eID uint64, p raylib.Vector2) bool {

	if qt.Boundry.Width == 0 || qt.Boundry.Height == 0 {
		utils.WarnLogger.Println("Failed to insert entity into QuadTree with zero width or height")
		return false
	}

	// If the position is not withing the boundry, reject the insertion
	if !qt.Boundry.Contains(p) {
		return false
	}

	// If there's still capacity, add the QuadTree hasn't been subdivided, add the entity
	if len(qt.Entities) < int(qt.Capacity) {
		qt.Entities = append(qt.Entities, eID)
		return true
	}

	// Subdivide if capacity is exceeded and not already divided
	if !qt.Divided {
		qt.Subdivide()
	}

	// Insert the entity into the appropriate child QuadTree
	if qt.NE.Insert(eID, p) {
		return true
	}
	if qt.NW.Insert(eID, p) {
		return true
	}
	if qt.SE.Insert(eID, p) {
		return true
	}
	if qt.SW.Insert(eID, p) {
		return true
	}

	// We should never reach this point
	utils.WarnLogger.Printf("Failed to insert entity %d into any child QuadTree", eID)
	return false
}

// Query retrieves all entities within a given range and appends them to 'found'
func (qt *QuadTree) Query(rangeRectt Rectangle, found *[]uint64) {

	// If the range does not intersect the boundry, return
	if !qt.Boundry.Intersects(&rangeRectt) {
		return
	}

	// If the range is completely within the boundry, add all entities
	for _, id := range qt.Entities {
		*found = append(*found, id)
	}

	// If subdivided, query the child QuadTrees
	if qt.Divided {
		qt.NE.Query(rangeRectt, found)
		qt.NW.Query(rangeRectt, found)
		qt.SE.Query(rangeRectt, found)
		qt.SW.Query(rangeRectt, found)
	}
}

// Clear removes all entities from the QuadTree and its children
func (qt *QuadTree) Clear() {

	qt.Entities = qt.Entities[:0]

	// Recursively clear all children
	if qt.Divided {

		// Run Clear on all children
		qt.NE.Clear()
		qt.NW.Clear()
		qt.SE.Clear()
		qt.SW.Clear()

		// Remove pointers to children
		qt.NE = nil
		qt.NW = nil
		qt.SE = nil
		qt.SW = nil

		// Set Divided to false
		qt.Divided = false
	}
}

/*

	// Future Implementation: Increment Update QuadTree

	// Steps to transition from rebuilding each frame to increment updates

	// Track Entity Positions: //
    // - Maintain a map of entity ID's to their previous positions.const
	// - Detect when an entity has moved by comparing the current and previous positions.

	// Implement Remove and Update functions:
	// - Remove: Delete an entity from the QuadTree based on its old position.const
	// - Update: Re-insert the entity with its new position if it has moved

	// Optimize QuadTree Structure:
	// - Add metadata for entities or maintain parent references to efficiently update the QuadTree.

	// Modify Collision System
	// - Instead of clearing and rebuilding the QuadTree each frame, only update entities that have moved.
	// - Reduce unnecessary insertions and deletions to enhance performance.

	// Handle Entity Destruction
	// - Ensure that entities removed from the game are also removed from the QuadTree to prevent stale refences

	// Concurrency Considerations
	// - If entities are updated concurrently, ensure that the QuadTree operations are thread-safe.

*/
