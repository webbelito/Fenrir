package ecs

import (
	"sort"
)

// Represents a system that can be updated, standard interface for all systems
type System interface {
	Update(dt float64, em *EntitiesManager, cm *ComponentsManager)
}

// Represents a system that can be rendered
type RenderableSystem interface {
	Render(*EntitiesManager, *ComponentsManager)
}

type SystemsManager struct {
	systems []SystemsWithPriority
}

type SystemsWithPriority struct {
	system   System
	priority int
}

func NewSystemsManager() *SystemsManager {
	return &SystemsManager{
		systems: []SystemsWithPriority{},
	}
}

func (sm *SystemsManager) AddSystem(system System, priority int) {
	sm.systems = append(sm.systems, SystemsWithPriority{system: system, priority: priority})

	// Sort systems by priority
	sort.SliceStable(sm.systems, func(i, j int) bool {
		return sm.systems[i].priority < sm.systems[j].priority
	})
}

func (sm *SystemsManager) Update(dt float64, em *EntitiesManager, cm *ComponentsManager) {
	for _, swp := range sm.systems {
		swp.system.Update(dt, em, cm)
	}
}

func (sm *SystemsManager) Render(em *EntitiesManager, cm *ComponentsManager) {

	// Render only on systems that is of type RenderableSystem
	for _, swp := range sm.systems {
		if rs, ok := swp.system.(RenderableSystem); ok {
			rs.Render(em, cm)
		}
	}
}
