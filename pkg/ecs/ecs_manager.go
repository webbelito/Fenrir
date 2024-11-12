package ecs

type ECSManager struct {
	entitiesManager   *EntitiesManager
	componentsManager *ComponentsManager
	systemsManager    *SystemsManager
}

func NewECSManager() *ECSManager {
	return &ECSManager{
		entitiesManager:   NewEntitiesManager(),
		componentsManager: NewComponentsManager(),
		systemsManager:    NewSystemsManager(),
	}
}

func (em *ECSManager) CreateEntity() Entity {
	return em.entitiesManager.CreateEntity()
}

func (em *ECSManager) DestroyEntity(entity Entity) {
	em.entitiesManager.DestroyEntity(entity)
	em.componentsManager.DestroyEntityComponents(entity)
}

func (em *ECSManager) GetAllEntities() []Entity {
	return em.entitiesManager.GetAllEntities()
}

func (em *ECSManager) GetEntityCount() int {
	return em.entitiesManager.GetEntityCount()
}

func (em *ECSManager) AddComponent(entity Entity, ct ComponentType, c Component) {
	em.componentsManager.AddComponent(entity, ct, c)
}

func (em *ECSManager) GetComponent(entity Entity, ct ComponentType) Component {
	return em.componentsManager.GetComponent(entity, ct)
}

func (em *ECSManager) AddSystem(system System, priority int) {
	em.systemsManager.AddSystem(system, priority)
}

func (em *ECSManager) Update(dt float64) {
	em.systemsManager.Update(dt, em.entitiesManager, em.componentsManager)
}
