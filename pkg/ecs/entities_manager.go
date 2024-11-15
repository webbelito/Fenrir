package ecs

import (
	"sync"
	"sync/atomic"
)

// Entity represents an unique entity with an unique ID
type Entity struct {
	ID uint64
}

// EntitiesManager manages the creation and destruction of entities
type EntitiesManager struct {
	nextEntityID uint64
	entities     map[uint64]*Entity
	entityMutex  sync.RWMutex
}

// NewEntitiesManager initializes and returns a new EntitiesManager
func NewEntitiesManager() *EntitiesManager {
	return &EntitiesManager{
		nextEntityID: 1,
		entities:     make(map[uint64]*Entity),
	}
}

// CreateEntity creates a new entity with a unique ID and adds it to the entities map
func (em *EntitiesManager) CreateEntity() *Entity {
	id := atomic.AddUint64(&em.nextEntityID, 1)
	entity := &Entity{ID: id}

	em.entityMutex.Lock()
	defer em.entityMutex.Unlock()

	em.entities[entity.ID] = entity
	return entity
}

// DestroyEntity removes an entity from the manager.
// It safely handles attempts to remove non-existent entities.
func (em *EntitiesManager) DestroyEntity(id uint64) {
	em.entityMutex.Lock()
	delete(em.entities, id)
	em.entityMutex.Unlock()
}

func (em *EntitiesManager) GetAllEntities() []*Entity {

	em.entityMutex.RLock()
	defer em.entityMutex.RUnlock()

	// Allocate a slice with the length of the entities map to avoid resizing
	entities := make([]*Entity, 0, len(em.entities))

	for entity := range em.entities {
		entities = append(entities, &Entity{ID: entity})
	}

	return entities
}

func (em *EntitiesManager) GetEntity(id uint64) (*Entity, bool) {
	em.entityMutex.RLock()
	defer em.entityMutex.RUnlock()

	entity, exists := em.entities[id]

	return entity, exists
}

func (em *EntitiesManager) GetEntityCount() int {
	em.entityMutex.RLock()
	defer em.entityMutex.RUnlock()

	return len(em.entities)
}

func (em *EntitiesManager) EntityExists(id uint64) bool {
	em.entityMutex.RLock()
	defer em.entityMutex.RUnlock()

	_, exists := em.entities[id]
	return exists
}
