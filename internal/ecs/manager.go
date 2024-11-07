// internal/ecs/manager.go

package ecs

import (
	"fmt"
	"sync"
)

type Manager struct {
	entities    map[Entity]map[string]Component
	systems     []System
	entityMutex sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		entities: make(map[Entity]map[string]Component),
		systems:  []System{},
	}
}

func (m *Manager) AddEntity(entity Entity) {
	m.entityMutex.Lock()
	defer m.entityMutex.Unlock()

	m.entities[entity] = make(map[string]Component)
}

func (m *Manager) RemoveEntity(entity Entity) {
	m.entityMutex.Lock()
	defer m.entityMutex.Unlock()

	delete(m.entities, entity)
}

func (m *Manager) AddComponent(entity Entity, component Component) {
	m.entityMutex.Lock()
	defer m.entityMutex.Unlock()

	m.entities[entity][m.getComponentName(component)] = component
}

func (m *Manager) GetComponent(entity Entity, componentType Component) Component {
	m.entityMutex.RLock()
	defer m.entityMutex.RUnlock()

	return m.entities[entity][m.getComponentName(componentType)]
}

func (m *Manager) AddSystem(system System) {
	m.systems = append(m.systems, system)
}

func (m *Manager) Update(deltaTime float32) {
	for _, system := range m.systems {
		system.Update(deltaTime, m)
	}
}

func (m *Manager) GetEntitiesWithComponents(component ...Component) []Entity {
	m.entityMutex.RLock()
	defer m.entityMutex.RUnlock()

	var result []Entity
	for entity, components := range m.entities {
		hasAllComponents := true
		for _, componentType := range component {
			if _, exists := components[m.getComponentName(componentType)]; !exists {
				hasAllComponents = false
				break
			}
		}

		if hasAllComponents {
			result = append(result, entity)
		}
	}

	return result
}

func (m *Manager) getComponentName(component Component) string {
	return fmt.Sprintf("%T", component)
}
