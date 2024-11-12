package ecs

import "sync"

type ComponentType uint64

const (
	PositionComponent ComponentType = iota
	VelocityComponent
	ColorComponent
	SpeedComponent
)

// Generic component interface
type Component interface{}

type ComponentsManager struct {
	Components map[ComponentType]map[Entity]Component
	compMutex  sync.RWMutex
}

func NewComponentsManager() *ComponentsManager {
	return &ComponentsManager{
		Components: make(map[ComponentType]map[Entity]Component),
	}
}

func (cm *ComponentsManager) AddComponent(e Entity, ct ComponentType, c Component) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	if _, exists := cm.Components[ct]; !exists {
		cm.Components[ct] = make(map[Entity]Component)
	}

	cm.Components[ct][e] = c
}

func (cm *ComponentsManager) GetComponent(e Entity, ct ComponentType) Component {

	cm.compMutex.RLock()
	defer cm.compMutex.RUnlock()

	if _, exists := cm.Components[ct]; exists {
		return cm.Components[ct][e]
	}

	return nil
}

func (cm *ComponentsManager) RemoveComponent(e Entity, ct ComponentType) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	if _, exists := cm.Components[ct]; exists {
		delete(cm.Components[ct], e)
	}
}

func (cm *ComponentsManager) DestroyEntityComponents(e Entity) {

	cm.compMutex.Lock()
	defer cm.compMutex.Unlock()

	for _, components := range cm.Components {
		delete(components, e)
	}
}
