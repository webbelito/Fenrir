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

func (em *ECSManager) CreateEntity() *Entity {
	return em.entitiesManager.CreateEntity()
}

func (em *ECSManager) DestroyEntity(id uint64) {
	em.entitiesManager.DestroyEntity(id)
	em.componentsManager.DestroyEntityComponents(id)
}

func (em *ECSManager) GetAllEntities() []*Entity {
	return em.entitiesManager.GetAllEntities()
}

func (em *ECSManager) GetEntityCount() int {
	return em.entitiesManager.GetEntityCount()
}

func (em *ECSManager) AddComponent(eID uint64, ct ComponentType, c Component) {
	em.componentsManager.AddComponent(eID, ct, c)
}

func (em *ECSManager) GetComponent(eID uint64, ct ComponentType) (Component, bool) {
	return em.componentsManager.GetComponent(eID, ct)
}

func (em *ECSManager) AddSystem(system System, priority int) {
	em.systemsManager.AddSystem(system, priority)
}

func (em *ECSManager) Update(dt float64) {
	em.systemsManager.Update(dt, em.entitiesManager, em.componentsManager)
}

func (em *ECSManager) Render() {
	em.systemsManager.Render(em.entitiesManager, em.componentsManager)
}
