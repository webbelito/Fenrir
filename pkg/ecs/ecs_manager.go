package ecs

import (
	"time"

	"github.com/webbelito/Fenrir/pkg/events"
	metricinterfaces "github.com/webbelito/Fenrir/pkg/interfaces/metricinterfaces"
	systeminterfaces "github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/utils"
)

type ECSManager struct {
	entitiesManager     *EntitiesManager
	componentsManager   *ComponentsManager
	uiComponentsManager *UIComponentsManager
	systemsManager      *SystemsManager
	eventsManager       *events.EventsManager

	performanceMetrics metricinterfaces.PerformanceMetrics
}

func NewECSManager() *ECSManager {

	// TODO: Implement a NewSystemsManager and NewPerformanceMetrics when creating the ECSManager
	ecsManager := &ECSManager{
		entitiesManager:     NewEntitiesManager(),
		componentsManager:   NewComponentsManager(),
		uiComponentsManager: NewUIComponentsManager(),
		eventsManager:       events.NewEventsManager(),
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
	em.uiComponentsManager.DestroyEntityComponents(id)
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

// * UIComponentsManager methods
func (em *ECSManager) GetUIComponentsManager() *UIComponentsManager {
	return em.uiComponentsManager
}

func (em *ECSManager) AddUIComponent(eID uint64, ct UIComponentType, c UIComponent) {
	em.uiComponentsManager.AddComponent(eID, ct, c)
}

func (em *ECSManager) GetUIComponent(eID uint64, ct UIComponentType) (UIComponent, bool) {
	return em.uiComponentsManager.GetComponent(eID, ct)
}

func (em *ECSManager) GetUIComponentsOfType(ct UIComponentType) (map[uint64]UIComponent, bool) {
	return em.uiComponentsManager.GetComponentsOfType(ct)
}

func (em *ECSManager) GetUIEntitiesWithComponents(cts []UIComponentType) []uint64 {
	return em.uiComponentsManager.GetEntitiesWithComponents(cts)
}

// * SystemsManager methods

func (em *ECSManager) GetSystemsManager() *SystemsManager {
	return em.systemsManager
}

func (em *ECSManager) AddLogicSystem(system systeminterfaces.UpdatableSystemInterface, priority int) {
	em.systemsManager.AddLogicSystem(system, priority)
}

func (em *ECSManager) RemoveLogicSystem(system systeminterfaces.UpdatableSystemInterface) {
	em.systemsManager.RemoveLogicSystem(system)
}

func (em *ECSManager) AddRenderSystem(system systeminterfaces.RenderableSystemInterface, priority int) {
	em.systemsManager.AddRenderSystem(system, priority)
}

func (em *ECSManager) RemoveRenderSystem(system systeminterfaces.RenderableSystemInterface) {
	em.systemsManager.RemoveRenderSystem(system)
}

func (em *ECSManager) AddUIRenderSystem(system systeminterfaces.UIRenderableSystemInterface, priority int) {
	em.systemsManager.AddUIRenderSystem(system, priority)
}

func (em *ECSManager) RemoveUIRenderSystem(system systeminterfaces.UIRenderableSystemInterface, priority int) {
	em.systemsManager.RemoveUIRenderSystem(system, priority)
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

func (em *ECSManager) RenderUISystems() {

	// TODO: Implement performance metrics for UI Render systems

	em.systemsManager.RenderUI()
}

// * EventsManager methods

func (em *ECSManager) GetEventsManager() *events.EventsManager {
	return em.eventsManager
}

// * PerformanceMetrics methods
func (em *ECSManager) GetPerformanceMetrics() metricinterfaces.PerformanceMetrics {
	return em.performanceMetrics
}

// * Player methods
// TODO: Currently hardcoded to a single player
func (ecsM *ECSManager) GetPlayerEntity() *Entity {

	// Find the entity with the player components
	playerComps := ecsM.componentsManager.GetEntitiesWithComponents([]ComponentType{PlayerComponent})

	playerID := uint64(0)
	for eID := range playerComps {

		if uint64(eID) == playerID {
			utils.WarnLogger.Println("Multiple entities with the player component found")
		}

		playerID = uint64(eID)
	}

	return nil
}

func (ecsM *ECSManager) GetCameraSystem() (systeminterfaces.CameraSystemInterface, bool) {
	return ecsM.systemsManager.GetCameraSystem()
}
