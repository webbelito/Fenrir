package ecs

import (
	"sync"
)

// UIComponentType represents the type of a UIComponent.
type UIComponentType uint64

// UIComponent is the interface that all UI components should implement
type UIComponent interface{}

// UIComponentsManager manages all UI components in the ECS
type UIComponentsManager struct {
	Components map[UIComponentType]map[uint64]UIComponent
	compMutex  sync.RWMutex
}

// NewUIComponentsManager initializes and returns a new UIComponentsManager
func NewUIComponentsManager() *UIComponentsManager {
	return &UIComponentsManager{
		Components: make(map[UIComponentType]map[uint64]UIComponent),
	}
}

// AddUIComponent adds a UI compnent to a specific entity
// It initializes the compnent map for the type if it doesn't exist
func (cm *UIComponentsManager) AddUIComponent(eID uint64, ct UIComponentType, c UIComponent) {
	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	// Initialize the component map if it doesn't exist
	if _, exists := cm.Components[ct]; !exists {
		cm.Components[ct] = make(map[uint64]UIComponent)
	}

	cm.Components[ct][eID] = c
}

// GetUIComponent retrieves a UI component of a specifc type for a given entity
// It return the compnent and a boolean indicating if the component exists.
func (cm *UIComponentsManager) GetUIComponent(eID uint64, ct UIComponentType) (UIComponent, bool) {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	// Get the component map for the type
	comps, compsExists := cm.Components[ct]
	if !compsExists {
		return nil, false
	}

	comp, compExists := comps[eID]
	return comp, compExists
}

// GetUIComponentsOfType retries all UI Components of a specific type.
// It returns a map and their components and a boolean indicating if the components exist.
func (cm *UIComponentsManager) GetUIComponentsOfType(ct UIComponentType) (map[uint64]UIComponent, bool) {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	// Get the component map for the type
	comps, compsExists := cm.Components[ct]
	if !compsExists {
		return nil, false
	}

	return comps, true
}

// GetUIEntitiesWithComponents retrieves all entity IDs that have all the specified UI component type.
func (cm *UIComponentsManager) GetUIEntitiesWithComponents(cts []UIComponentType) []uint64 {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	// If no component types are provided, return nil
	if len(cts) == 0 {
		return nil
	}

	// Initialize with the smallest compnent type set to optimize the search.
	minIndex := 0
	minCount := len(cm.Components[cts[0]])

	// Find the smallest component type set
	for i, ct := range cts {

		// If a component type is found, check if it has the least number of components
		if comps, compsExists := cm.Components[ct]; compsExists {

			// If the component type has less components than the current smallest, update the smallest
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

	// Iterate over the smallest component type set
	for entity := range smallestComps {
		hasAll := true

		// Check if the entity has all the other component types
		for i, ct := range cts {

			// Skip the smallest component type
			if i == minIndex {
				continue
			}

			// If the entity does not have a component of the type, break the loop
			if _, compExists := cm.Components[ct][entity]; !compExists {
				hasAll = false
				break
			}
		}

		// If the entity has all the component types, add it to the list
		if hasAll {
			entities = append(entities, entity)
		}
	}

	return entities
}

// RemoveUIComponent removes a UI component of a specific type from an entity.
func (cm *UIComponentsManager) RemoveUIComponent(eID uint64, ct UIComponentType) {
	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	// Get the component map for the type
	if comps, compsExists := cm.Components[ct]; compsExists {
		delete(comps, eID)
	}
}

// DestroyUIEntityComponents removes all UI components associated with an entity
func (cm *UIComponentsManager) DestroyUIEntityComponents(eID uint64) {
	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	// Iterate over all component types
	for _, components := range cm.Components {
		delete(components, eID)
	}
}
