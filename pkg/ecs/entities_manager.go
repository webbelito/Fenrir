package ecs

type Entity struct {
	ID uint64
}

type EntitiesManager struct {
	nextEntityID uint64
	entites      map[uint64]bool
}

func NewEntitiesManager() *EntitiesManager {
	return &EntitiesManager{
		nextEntityID: 1,
		entites:      make(map[uint64]bool),
	}
}

func (em *EntitiesManager) CreateEntity() *Entity {
	entity := &Entity{ID: em.nextEntityID}
	em.entites[entity.ID] = true
	em.nextEntityID++

	return entity
}

func (em *EntitiesManager) DestroyEntity(e *Entity) {
	delete(em.entites, e.ID)
}

func (em *EntitiesManager) GetAllEntities() []*Entity {
	var entities []*Entity

	for entity := range em.entites {
		entities = append(entities, &Entity{ID: entity})
	}

	return entities
}

func (em *EntitiesManager) GetEntity(id uint64) *Entity {
	if _, exists := em.entites[id]; exists {
		return &Entity{ID: id}
	}

	return nil
}

func (em *EntitiesManager) GetEntityCount() int {
	return len(em.entites)
}
