package ecs

import (
	"time"

	metricinterfaces "github.com/webbelito/Fenrir/pkg/interfaces/metricinterfaces"
	systeminterfaces "github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
)

type ECSManager struct {
	entitiesManager   *EntitiesManager
	componentsManager *ComponentsManager
	systemsManager    *SystemsManager

	performanceMetrics metricinterfaces.PerformanceMetrics
}

func NewECSManager() *ECSManager {

	ecsManager := &ECSManager{
		entitiesManager:   NewEntitiesManager(),
		componentsManager: NewComponentsManager(),
	}

	ecsManager.systemsManager = NewSystemsManager(ecsManager)

	ecsManager.performanceMetrics = metricinterfaces.PerformanceMetrics{}

	return ecsManager
}

// * EntitiesManager methods

func (em *ECSManager) GetEntitiesManager() *EntitiesManager {
	return em.entitiesManager
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

// * ComponentsManager methods

func (em *ECSManager) GetComponentsManager() *ComponentsManager {
	return em.componentsManager
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

// * SystemsManager methods

func (em *ECSManager) GetSystemsManager() *SystemsManager {
	return em.systemsManager
}

func (em *ECSManager) AddLogicSystem(system systeminterfaces.Updatable, priority int) {
	em.systemsManager.AddLogicSystem(system, priority)
}

func (em *ECSManager) RemoveLogicSystem(system systeminterfaces.Updatable) {
	em.systemsManager.RemoveLogicSystem(system)
}

func (em *ECSManager) AddRenderSystem(system systeminterfaces.Renderable, priority int) {
	em.systemsManager.AddRenderSystem(system, priority)
}

func (em *ECSManager) RemoveRenderSystem(system systeminterfaces.Renderable) {
	em.systemsManager.RemoveRenderSystem(system)
}

func (em *ECSManager) UpdateLogicSystems(dt float64) {
	em.performanceMetrics.UpdateStartTime = time.Now()
	em.systemsManager.Update(dt)
	em.performanceMetrics.UpdateDuration = time.Since(em.performanceMetrics.UpdateStartTime)
}

func (em *ECSManager) UpdateRenderSystems() {
	em.performanceMetrics.RenderStartTime = time.Now()
	em.systemsManager.Render()
	em.performanceMetrics.RenderDuration = time.Since(em.performanceMetrics.RenderStartTime)

	// Update the total time
	em.performanceMetrics.TotalDuration = time.Since(em.performanceMetrics.UpdateStartTime)
}

// * PerformanceMetrics methods
func (em *ECSManager) GetPerformanceMetrics() metricinterfaces.PerformanceMetrics {
	return em.performanceMetrics
}
