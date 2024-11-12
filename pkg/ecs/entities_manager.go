package ecs

type Entity uint64

type EntitiesManager struct {
	nextEntityID Entity
	entites      map[Entity]bool
}

func NewEntitiesManager() *EntitiesManager {
	return &EntitiesManager{
		nextEntityID: 1,
		entites:      make(map[Entity]bool),
	}
}

func (em *EntitiesManager) CreateEntity() Entity {
	entity := em.nextEntityID
	em.entites[entity] = true
	em.nextEntityID++

	return entity
}

func (em *EntitiesManager) DestroyEntity(entity Entity) {
	delete(em.entites, entity)
}

func (em *EntitiesManager) GetAllEntities() []Entity {
	var entities []Entity

	for entity := range em.entites {
		entities = append(entities, entity)
	}

	return entities
}

func (em *EntitiesManager) GetEntityCount() int {
	return len(em.entites)
}
