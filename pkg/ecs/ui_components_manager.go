package ecs

import (
	"sync"
)

// UIComponentType represents the type of a UIComponent.
type UIComponentType uint64

// UIComponent is the interface that all UI components should implement
type UIComponent interface{}

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

// AddComponent adds a UI compnent to a specific entity
// It initializes the compnent map for the type if it doesn't exist
func (cm *UIComponentsManager) AddComponent(eID uint64, ct UIComponentType, c UIComponent) {
	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	if _, exists := cm.Components[ct]; !exists {
		cm.Components[ct] = make(map[uint64]UIComponent)
	}

	cm.Components[ct][eID] = c
}

// GetComponent retrieves a UI component of a specifc type for a given entity
// It return the compnent and a boolean indicating if the component exists.
func (cm *UIComponentsManager) GetComponent(eID uint64, ct UIComponentType) (UIComponent, bool) {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	comps, compsExists := cm.Components[ct]
	if !compsExists {
		return nil, false
	}

	comp, compExists := comps[eID]
	return comp, compExists
}

// GetComponentsOfType retries all UI Components of a specific type.
// It returns a map and their components and a boolean indicating if the components exist.
func (cm *UIComponentsManager) GetComponentsOfType(ct UIComponentType) (map[uint64]UIComponent, bool) {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	comps, compsExists := cm.Components[ct]
	if !compsExists {
		return nil, false
	}

	return comps, true
}

// GetEntitiesWithComponents retrieves all entity IDs that have all the specified UI component type.
func (cm *UIComponentsManager) GetEntitiesWithComponents(cts []UIComponentType) []uint64 {
	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	if len(cts) == 0 {
		return nil
	}

	// Initialize with the smallest compnent type set to optimize the search.
	minIndex := 0
	minCount := len(cm.Components[cts[0]])

	for i, ct := range cts {
		if comps, compsExists := cm.Components[ct]; compsExists {
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
			if _, compExists := cm.Components[ct][entity]; !compExists {
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

// RemoveComponent removes a UI component of a specific type from an entity.
func (cm *UIComponentsManager) RemoveComponent(eID uint64, ct UIComponentType) {
	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	if comps, compsExists := cm.Components[ct]; compsExists {
		delete(comps, eID)
	}
}

// DestroyEntityComponents removes all UI components associated with an entity
func (cm *UIComponentsManager) DestroyEntityComponents(eID uint64) {
	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	for _, components := range cm.Components {
		delete(components, eID)
	}
}
