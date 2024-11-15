package ecs

import (
	"sync"
)

type ComponentType uint64

// Generic component interface
type Component interface{}

type ComponentsManager struct {
	Components map[ComponentType]map[*Entity]Component
	compMutex  sync.RWMutex
}

func NewComponentsManager() *ComponentsManager {
	return &ComponentsManager{
		Components: make(map[ComponentType]map[*Entity]Component),
	}
}

func (cm *ComponentsManager) AddComponent(e *Entity, ct ComponentType, c Component) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	if _, exists := cm.Components[ct]; !exists {
		cm.Components[ct] = make(map[*Entity]Component)
	}

	cm.Components[ct][e] = c
}

func (cm *ComponentsManager) GetComponent(e *Entity, ct ComponentType) Component {

	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	if _, exists := cm.Components[ct]; exists {
		return cm.Components[ct][e]
	}

	return nil
}

func (cm *ComponentsManager) GetEntitiesWithComponents(componentTypes []ComponentType) []*Entity {

	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	if len(componentTypes) == 0 {
		return nil
	}

	// Find the component type with the least amount of entities
	minCompType := componentTypes[0]
	minCount := len(cm.Components[minCompType])

	for _, ct := range componentTypes {

		entitiesMap, exists := cm.Components[ct]
		if !exists || len(entitiesMap) == 0 {
			return nil
		}

		if len(entitiesMap) < minCount {
			minCount = len(entitiesMap)
			minCompType = ct
		}
	}

	// Start with entities having the least common component
	baseEntities := cm.Components[minCompType]
	result := make([]*Entity, 0)

	// Check each entity for the presence of other components
	for entity := range baseEntities {
		hasAllComponents := true

		// Iterate over all component types
		for _, ct := range componentTypes {
			if ct == minCompType {
				continue
			}
			if _, exists := cm.Components[ct][entity]; !exists {
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

func (cm *ComponentsManager) RemoveComponent(e *Entity, ct ComponentType) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	if _, exists := cm.Components[ct]; exists {
		delete(cm.Components[ct], e)
	}
}

func (cm *ComponentsManager) DestroyEntityComponents(e *Entity) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	for _, components := range cm.Components {
		delete(components, e)
	}
}
