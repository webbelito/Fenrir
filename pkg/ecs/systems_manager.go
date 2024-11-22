package ecs

import (
	"sort"

	systeminterfaces "github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/utils"
)

type SystemsManager struct {
	logicSystems  []systeminterfaces.Updatable
	renderSystems []systeminterfaces.Renderable
}

func NewSystemsManager(ecsManager *ECSManager) *SystemsManager {
	return &SystemsManager{
		logicSystems:  []systeminterfaces.Updatable{},
		renderSystems: []systeminterfaces.Renderable{},
	}
}

func (sm *SystemsManager) AddLogicSystem(system systeminterfaces.Updatable, priority int) {

	inserted := false

	// TODO: Implement priority

	if !inserted {
		sm.logicSystems = append(sm.logicSystems, system)

		sort.Slice(sm.logicSystems, func(i, j int) bool {
			return sm.logicSystems[i].GetPriority() < sm.logicSystems[j].GetPriority()
		})
	}

	utils.InfoLogger.Printf("Added logic system: %T\n", system)

}

func (sm *SystemsManager) RemoveLogicSystem(system systeminterfaces.Updatable) {
	for i, sys := range sm.logicSystems {
		if sys == system {
			sm.logicSystems = append(sm.logicSystems[:i], sm.logicSystems[i+1:]...)
			utils.InfoLogger.Printf("Removed logic system: %T\n", system)
			break
		}
	}
}

func (sm *SystemsManager) AddRenderSystem(system systeminterfaces.Renderable, priority int) {

	inserted := false

	// TODO: Implement priority

	if !inserted {
		sm.renderSystems = append(sm.renderSystems, system)

		sort.Slice(sm.renderSystems, func(i, j int) bool {
			return sm.renderSystems[i].GetPriority() < sm.renderSystems[j].GetPriority()
		})
	}

	utils.InfoLogger.Printf("Added render system: %T\n", system)

}

func (sm *SystemsManager) RemoveRenderSystem(system systeminterfaces.Renderable) {
	for i, sys := range sm.renderSystems {
		if sys == system {
			sm.renderSystems = append(sm.renderSystems[:i], sm.renderSystems[i+1:]...)
			utils.InfoLogger.Printf("Removed render system: %T\n", system)
			break
		}
	}
}

// GetCameraSystem retrieves the CameraSystem from logic or render systems.
func (sm *SystemsManager) GetCameraSystem() (systeminterfaces.CameraSystemInterface, bool) {
	// Check Logic Systems
	for _, sys := range sm.logicSystems {
		if cameraSys, ok := sys.(systeminterfaces.CameraSystemInterface); ok {
			return cameraSys, true
		}
	}

	// Check Render Systems
	for _, sys := range sm.renderSystems {
		if cameraSys, ok := sys.(systeminterfaces.CameraSystemInterface); ok {
			return cameraSys, true
		}
	}

	return nil, false
}

func (sm *SystemsManager) Update(dt float64) {
	for _, system := range sm.logicSystems {
		system.Update(dt)
	}
}

func (sm *SystemsManager) Render() {
	for _, system := range sm.renderSystems {
		system.Render()
	}
}
