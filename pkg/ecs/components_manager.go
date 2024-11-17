package ecs

import (
	"sync"
)

// ComponentType represent the type of a component
type ComponentType uint64

// Component is the interface that all components should implement
type Component interface{}

// ComponentsManager manages all components associated with entities
type ComponentsManager struct {

	// TODO: Rename Components to components, ensure that the map is private
	Components map[ComponentType]map[uint64]Component
	compMutex  sync.RWMutex
}

// NewComponentsManager initliazes and returns a new ComponentsManager
func NewComponentsManager() *ComponentsManager {
	return &ComponentsManager{
		Components: make(map[ComponentType]map[uint64]Component),
	}
}

// AddComponent adds a component to a specific entity.
// It initializes the compopnent map for the type if it doesn't exist.
func (cm *ComponentsManager) AddComponent(eID uint64, ct ComponentType, c Component) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	if _, exists := cm.Components[ct]; !exists {
		cm.Components[ct] = make(map[uint64]Component)
	}

	cm.Components[ct][eID] = c
}

// GetComponent retrieves a component of a specific type for a given entity.
// It returns the compenent and a boolean indicating if the component exists.
func (cm *ComponentsManager) GetComponent(eID uint64, ct ComponentType) (Component, bool) {

	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	comps, compsExists := cm.Components[ct]

	if !compsExists {
		return nil, false
	}

	comp, compExists := comps[eID]
	return comp, compExists

}

// GetComponentsOfType retrieves all components of a specific type.
// It returns a map of entities and their components and a boolean indicating if the components exist.
func (cm *ComponentsManager) GetComponentsOfType(ct ComponentType) (map[uint64]Component, bool) {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	comps, compsExists := cm.Components[ct]

	if !compsExists {
		return nil, false
	}

	return comps, compsExists
}

func (cm *ComponentsManager) GetEntitiesWithComponents(cts []ComponentType) []uint64 {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	if len(cts) == 0 {
		return nil
	}

	// Initialize with the smallest component type set to optimize the search.
	minIndex := 0
	minCount := len(cm.Components[cts[0]])

	for i, ct := range cts {
		if comps, exists := cm.Components[ct]; exists {
			if len(comps) < minCount {
				minCount = len(comps)
				minIndex = i
			}
		} else {
			// If a component type is not found, return nil
			return nil
		}
	}

	// Start with entities having the least common component
	smallestComps := cm.Components[cts[minIndex]]
	entities := make([]uint64, 0, minCount)

	for entity := range smallestComps {
		hasAll := true
		for i, ct := range cts {
			if i == minIndex {
				continue
			}
			if _, exists := cm.Components[ct][entity]; !exists {
				hasAll = false
				break
			}
		}
		if hasAll {
			entities = append(entities, entity)
		}
	}

	return entities

}

func (cm *ComponentsManager) RemoveComponent(eID uint64, ct ComponentType) {

	cm.compMutex.Lock()
	delete(cm.Components[ct], eID)
	defer cm.compMutex.Unlock()

}

func (cm *ComponentsManager) DestroyEntityComponents(id uint64) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	for _, components := range cm.Components {
		delete(components, id)
	}
}
