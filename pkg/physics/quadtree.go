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
	Boundry      Rectangle
	Capacity     int32
	Entities     []uint64
	Divided      bool
	MaxDepth     int32
	CurrentDepth int32
	NE           *QuadTree
	NW           *QuadTree
	SE           *QuadTree
	SW           *QuadTree
}

// NewQuadTree initializes a new QuadTree node
func NewQuadTree(b Rectangle, c int32, mD int32, cD int32) *QuadTree {
	return &QuadTree{
		Boundry:      b,
		Capacity:     c,
		Entities:     make([]uint64, 0, c),
		Divided:      false,
		MaxDepth:     mD,
		CurrentDepth: cD,
	}
}

// Subdivide splits the QuadTree into four childen QuadTrees
func (qt *QuadTree) Subdivide() {
	halfWidth := qt.Boundry.Width / 2
	halfHeight := qt.Boundry.Height / 2
	x := qt.Boundry.Position.X
	y := qt.Boundry.Position.Y

	qt.NE = NewQuadTree(Rectangle{raylib.NewVector2(x+halfWidth, y), halfWidth, halfHeight}, qt.Capacity, qt.MaxDepth, qt.CurrentDepth+1)
	qt.NW = NewQuadTree(Rectangle{raylib.NewVector2(x, y), halfWidth, halfHeight}, qt.Capacity, qt.MaxDepth, qt.CurrentDepth+1)
	qt.SE = NewQuadTree(Rectangle{raylib.NewVector2(x+halfWidth, y+halfHeight), halfWidth, halfHeight}, qt.Capacity, qt.MaxDepth, qt.CurrentDepth+1)
	qt.SW = NewQuadTree(Rectangle{raylib.NewVector2(x, y+halfHeight), halfWidth, halfHeight}, qt.Capacity, qt.MaxDepth, qt.CurrentDepth+1)

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
		if qt.CurrentDepth >= qt.MaxDepth {
			// We've reached the maximum depth, do not subdivide
			qt.Entities = append(qt.Entities, eID)
			return true
		}
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
