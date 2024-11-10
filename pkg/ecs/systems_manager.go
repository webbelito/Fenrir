package ecs

import "sort"

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
