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

func (em *ECSManager) DestroyEntity(e *Entity) {
	em.entitiesManager.DestroyEntity(e)
	em.componentsManager.DestroyEntityComponents(e)
}

func (em *ECSManager) GetAllEntities() []*Entity {
	return em.entitiesManager.GetAllEntities()
}

func (em *ECSManager) GetEntityCount() int {
	return em.entitiesManager.GetEntityCount()
}

func (em *ECSManager) AddComponent(e *Entity, ct ComponentType, c Component) {
	em.componentsManager.AddComponent(e, ct, c)
}

func (em *ECSManager) GetComponent(e *Entity, ct ComponentType) Component {
	return em.componentsManager.GetComponent(e, ct)
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
