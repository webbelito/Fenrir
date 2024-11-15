package ecs

type Entity struct {
	ID uint64
}

type EntitiesManager struct {
	nextEntityID uint64
	entities     map[uint64]*Entity
}

func NewEntitiesManager() *EntitiesManager {
	return &EntitiesManager{
		nextEntityID: 1,
		entities:     make(map[uint64]*Entity),
	}
}

func (em *EntitiesManager) CreateEntity() *Entity {
	entity := &Entity{ID: em.nextEntityID}
	em.entities[entity.ID] = entity
	em.nextEntityID++

	return entity
}

func (em *EntitiesManager) DestroyEntity(e *Entity) {
	delete(em.entities, e.ID)
}

func (em *EntitiesManager) GetAllEntities() []*Entity {
	entities := make([]*Entity, 0, len(em.entities))

	for entity := range em.entities {
		entities = append(entities, &Entity{ID: entity})
	}

	return entities
}

func (em *EntitiesManager) GetEntity(id uint64) (*Entity, bool) {
	entity, exists := em.entities[id]
	return entity, exists
}

func (em *EntitiesManager) GetEntityCount() int {
	return len(em.entities)
}
