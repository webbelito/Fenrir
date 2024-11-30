package ecs

import (
	"errors"
	"sync"

	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/events"
	metricinterfaces "github.com/webbelito/Fenrir/pkg/interfaces/metricinterfaces"
	"github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/resources"
	"github.com/webbelito/Fenrir/pkg/utils"
)

// Manager is the main struct that holds all the managers and systems
type Manager struct {
	entities     *EntitiesManager
	components   *ComponentsManager
	uiComponents *UIComponentsManager
	systems      *SystemsManager
	events       *events.EventsManager
	resources    *resources.ResourcesManager
	metricinterfaces.PerformanceMetrics
	mutex sync.Mutex
}

// NewManager creates a new Manager
func NewManager() *Manager {

	manager := &Manager{
		entities:     NewEntitiesManager(),
		components:   NewComponentsManager(),
		uiComponents: NewUIComponentsManager(),
		systems:      NewSystemsManager(),
		events:       events.NewEventsManager(),
		resources:    resources.NewResourcesManager(),
	}

	manager.PerformanceMetrics = metricinterfaces.PerformanceMetrics{}

	return manager
}

// * Entities

// CreateEntity creates a new entity
func (m *Manager) CreateEntity() *Entity {
	return m.entities.CreateEntity()
}

// DestroyEntity destroys an entity and its related components
func (m *Manager) DestroyEntity(eID uint64) {
	m.entities.DestroyEntity(eID)
	m.components.DestroyEntityComponents(eID)
	m.uiComponents.DestroyEntityComponents(eID)
}

// GetAllEntities returns all entities
func (m *Manager) GetAllEntities() []*Entity {
	return m.entities.GetAllEntities()
}

// GetEntityCount returns the number of entities
func (m *Manager) GetEntityCount() int {
	return m.entities.GetEntityCount()
}

// * Components

/*
GetComponet returns a component of a given type for an entity
Component is an empty interface, so it can be any type. Remember to type assert it when using it.
*/

func (m *Manager) GetComponent(eID uint64, ct ComponentType) (Component, bool) {
	return m.components.GetComponent(eID, ct)
}

// GetComponentsOfType retrieves all components of a specific type.
// It returns a map of entity IDs to their components and a boolean indicating existence.
func (m *Manager) GetComponentsOfType(ct ComponentType) (map[uint64]Component, bool) {
	return m.components.GetComponentsOfType(ct)
}

// GetEntitiesWithComponents retrieves all entity IDs that have all the specified component types.
func (m *Manager) GetEntitiesWithComponents(cts []ComponentType) []uint64 {
	return m.components.GetEntitiesWithComponents(cts)
}

// AddComponent adds a component to an entity
func (m *Manager) AddComponent(eID uint64, ct ComponentType, c Component) {
	m.components.AddComponent(eID, ct, c)
}

// RemoveComponent removes a component from an entity
func (m *Manager) RemoveComponent(eID uint64, ct ComponentType) {
	m.components.RemoveComponent(eID, ct)
}

// GetCamera retrieves the Camera component.
func (m *Manager) GetCamera() (*components.Camera, bool) {
	cameraComp, exists := m.GetComponent(0, CameraComponent) // Assuming entity ID 0 is the camera
	if !exists {
		return nil, false
	}
	camera, ok := cameraComp.(*components.Camera)
	return camera, ok
}

// * UI Components

// GetUIComponent returns a UI component of a given type for an entity
func (m *Manager) GetUIComponent(eID uint64, ct UIComponentType) (UIComponent, bool) {
	return m.uiComponents.GetComponent(eID, ct)
}

// GetComponentsOfType retrieves all UI components of a specific type.
// It returns a map of entity IDs to their components and a boolean indicating existence.
func (m *Manager) GetUIComponentsOfType(ct UIComponentType) (map[uint64]UIComponent, bool) {
	return m.uiComponents.GetComponentsOfType(ct)
}

// GetUIEntitiesWithComponents retrieves all entity IDs that have all the specified UI component types.
func (m *Manager) GetUIEntitiesWithComponents(cts []UIComponentType) []uint64 {
	return m.uiComponents.GetEntitiesWithComponents(cts)
}

// AddUIComponent adds a UI component to an entity
func (m *Manager) AddUIComponent(eID uint64, ct UIComponentType, c UIComponent) {
	m.uiComponents.AddComponent(eID, ct, c)
}

// RemoveUIComponent removes a UI component from an entity
func (m *Manager) RemoveUIComponent(eID uint64, ct UIComponentType) {
	m.uiComponents.RemoveComponent(eID, ct)
}

// * Systems

// RegisterSystem registers a system based on its type, returns an error if the systemtype is not found
func (m *Manager) RegisterSystem(s interface{}, p int) error {
	switch system := s.(type) {
	case systeminterfaces.UpdatableSystemInterface:
		m.systems.AddLogicSystem(system, p)
	case systeminterfaces.RenderableSystemInterface:
		m.systems.AddRenderSystem(system, p)
	case systeminterfaces.UIRenderableSystemInterface:
		m.systems.AddUIRenderSystem(system, p)
	default:
		utils.ErrorLogger.Println("Unknown system type: ", system)
		return errors.New("unknown system type")
	}

	return nil

}

// UnregisterSystem unregisters a system based on its type, returns an error if the systemtype is not found
func (m *Manager) UnregisterSystem(s interface{}) error {
	switch system := s.(type) {
	case systeminterfaces.UpdatableSystemInterface:
		m.systems.RemoveLogicSystem(system)
	case systeminterfaces.RenderableSystemInterface:
		m.systems.RemoveRenderSystem(system)
	case systeminterfaces.UIRenderableSystemInterface:
		m.systems.RemoveUIRenderSystem(system)
	default:
		utils.ErrorLogger.Println("Unknown system type: ", system)
		return errors.New("unknown system type")
	}

	return nil
}

// Update updates all systems
func (m *Manager) Update(dt float64) {

	// Thread safetey
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.systems.Update(dt)
}

// Render renders all systems, including UI
func (m *Manager) Render() {

	// Thread safetey
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.systems.Render()
	m.systems.RenderUI()
}

// * Resources

// LoadTexture loads a texture from a file and stores it in the ResourceManager
func (m *Manager) LoadTexture(path string) (raylib.Texture2D, error) {
	return m.resources.LoadTexture(path)
}

// GetTexture gets a texture from the ResourceManager
func (m *Manager) GetTexture(path string) (raylib.Texture2D, bool) {
	return m.resources.GetTexture(path)
}

// UnloadTextures unloads all textures from the ResourceManager
func (m *Manager) UnloadTextures() {
	m.resources.UnloadTextures()
}

// LoadSound loads a sound from a file and stores it in the ResourceManager
func (m *Manager) LoadSound(path string) (raylib.Sound, error) {
	return m.resources.LoadSound(path)
}

// * Events

// RegisterEvent registers an event with the EventsManager
func (m *Manager) DispatchEvent(eventName string, event interface{}) {
	m.events.Dispatch(eventName, event)
}

// Subscribe subscribes a system to an event
func (m *Manager) SubscribeEvent(eventName string, system events.EventsHandler) {
	m.events.Subscribe(eventName, system)
}

// * Metrics

// GetPerformanceMetrics returns the performance metrics
func (m *Manager) GetPerformanceMetrics() metricinterfaces.PerformanceMetrics {
	return m.PerformanceMetrics
}
