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

	// Initialize the component map if it doesn't exist
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

	// Get the component map for the type
	comps, compsExists := cm.Components[ct]

	// If the component map doesn't exist, return false
	if !compsExists {
		return nil, false
	}

	// Get the component from the map
	comp, compExists := comps[eID]
	return comp, compExists

}

// GetComponentsOfType retrieves all components of a specific type.
// It returns a map of entities and their components and a boolean indicating if the components exist.
func (cm *ComponentsManager) GetComponentsOfType(ct ComponentType) (map[uint64]Component, bool) {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	// Get the component map for the type
	comps, compsExists := cm.Components[ct]

	// If the component map doesn't exist, return false
	if !compsExists {
		return nil, false
	}

	return comps, compsExists
}

func (cm *ComponentsManager) GetEntitiesWithComponents(cts []ComponentType) []uint64 {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	// If no component types are provided, return nil
	if len(cts) == 0 {
		return nil
	}

	// Initialize with the smallest component type set to optimize the search.
	minIndex := 0
	minCount := len(cm.Components[cts[0]])

	// Find the component type with the least number of components
	for i, ct := range cts {

		// If a component type is found, check if it has the least number of components
		if comps, exists := cm.Components[ct]; exists {

			// If the number of components is less than the current minimum, update the minimum
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

	// Iterate over the entities with the smallest component type
	for entity := range smallestComps {
		hasAll := true

		// Check if the entity has all the required components
		for i, ct := range cts {

			// Skip the smallest component type
			if i == minIndex {
				continue
			}

			// If the entity doesn't have a component of the type, break the loop
			if _, exists := cm.Components[ct][entity]; !exists {
				hasAll = false
				break
			}
		}

		// If the entity has all the required components, add it to the list
		if hasAll {
			entities = append(entities, entity)
		}
	}

	return entities

}

func (cm *ComponentsManager) RemoveComponent(eID uint64, ct ComponentType) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	// Remove the component from the map
	delete(cm.Components[ct], eID)

}

func (cm *ComponentsManager) DestroyEntityComponents(id uint64) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	// Remove the entity's components from all component maps
	for _, components := range cm.Components {
		delete(components, id)
	}
}
