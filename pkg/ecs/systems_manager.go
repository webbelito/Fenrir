package ecs

import (
	"sync"

	systeminterfaces "github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/utils"
)

type SystemsManager struct {
	logicSystems    []systeminterfaces.UpdatableSystemInterface
	renderSystems   []systeminterfaces.RenderableSystemInterface
	uiRenderSystems []systeminterfaces.UIRenderableSystemInterface
	systemMutex     sync.RWMutex
}

func NewSystemsManager(ecsManager *ECSManager) *SystemsManager {
	return &SystemsManager{
		logicSystems:    []systeminterfaces.UpdatableSystemInterface{},
		renderSystems:   []systeminterfaces.RenderableSystemInterface{},
		uiRenderSystems: []systeminterfaces.UIRenderableSystemInterface{},
	}
}

func (sm *SystemsManager) AddLogicSystem(system systeminterfaces.UpdatableSystemInterface, priority int) {

	//sm.systemMutex.Lock()
	//defer sm.systemMutex.Unlock()

	inserted := false

	// TODO: Implement priority

	for i, existingSystem := range sm.logicSystems {
		if existingSystem.GetPriority() > system.GetPriority() {
			sm.logicSystems = append(sm.logicSystems[:i], append([]systeminterfaces.UpdatableSystemInterface{system}, sm.logicSystems[i:]...)...)

			inserted = true

			utils.InfoLogger.Printf("Added logic system: %T\n", system)

			break
		}
	}

	if !inserted {
		sm.logicSystems = append(sm.logicSystems, system)
		utils.InfoLogger.Printf("Added logic system: %T\n", system)
	}

}

func (sm *SystemsManager) RemoveLogicSystem(system systeminterfaces.UpdatableSystemInterface) {

	sm.systemMutex.Lock()
	defer sm.systemMutex.Unlock()

	for i, sys := range sm.logicSystems {
		if sys == system {
			sm.logicSystems = append(sm.logicSystems[:i], sm.logicSystems[i+1:]...)
			utils.InfoLogger.Printf("Removed logic system: %T\n", system)
			break
		}
	}
}

func (sm *SystemsManager) AddRenderSystem(system systeminterfaces.RenderableSystemInterface, priority int) {

	//sm.systemMutex.Lock()
	//defer sm.systemMutex.Unlock()

	inserted := false

	// TODO: Implement priority

	for i, existingSystem := range sm.renderSystems {
		if existingSystem.GetPriority() > system.GetPriority() {
			sm.renderSystems = append(sm.renderSystems[:i], append([]systeminterfaces.RenderableSystemInterface{system}, sm.renderSystems[i:]...)...)
			inserted = true
			utils.InfoLogger.Printf("Added render system: %T\n", system)
			break
		}
	}

	if !inserted {
		sm.renderSystems = append(sm.renderSystems, system)
		utils.InfoLogger.Printf("Added render system: %T\n", system)
	}

}

func (sm *SystemsManager) RemoveRenderSystem(system systeminterfaces.RenderableSystemInterface) {

	sm.systemMutex.Lock()
	defer sm.systemMutex.Unlock()

	for i, sys := range sm.renderSystems {
		if sys == system {
			sm.renderSystems = append(sm.renderSystems[:i], sm.renderSystems[i+1:]...)
			utils.InfoLogger.Printf("Removed render system: %T\n", system)
			break
		}
	}
}

func (sm *SystemsManager) AddUIRenderSystem(system systeminterfaces.UIRenderableSystemInterface, priority int) {

	sm.systemMutex.Lock()
	defer sm.systemMutex.Unlock()

	inserted := false

	for i, existingSystem := range sm.uiRenderSystems {
		if existingSystem.GetPriority() > system.GetPriority() {
			sm.uiRenderSystems = append(sm.uiRenderSystems[:i], append([]systeminterfaces.UIRenderableSystemInterface{system}, sm.uiRenderSystems[i:]...)...)
			inserted = true
			utils.InfoLogger.Printf("Added UI render system: %T\n", system)
			break
		}
	}

	if !inserted {
		sm.uiRenderSystems = append(sm.uiRenderSystems, system)
	}
}

func (sm *SystemsManager) RemoveUIRenderSystem(system systeminterfaces.UIRenderableSystemInterface, priority int) {

	sm.systemMutex.Lock()
	defer sm.systemMutex.Unlock()

	for i, sys := range sm.uiRenderSystems {
		if sys == system {
			sm.uiRenderSystems = append(sm.uiRenderSystems[:i], sm.uiRenderSystems[i+1:]...)
			utils.InfoLogger.Printf("Removed UI render system: %T\n", system)
			break
		}
	}
}

// GetCameraSystem retrieves the CameraSystem from logic or render systems.
func (sm *SystemsManager) GetCameraSystem() (systeminterfaces.CameraSystemInterface, bool) {

	sm.systemMutex.RLock()
	defer sm.systemMutex.RUnlock()

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

	sm.systemMutex.RLock()
	defer sm.systemMutex.RUnlock()

	for _, system := range sm.logicSystems {
		system.Update(dt)
	}
}

func (sm *SystemsManager) Render() {
	sm.systemMutex.RLock()
	defer sm.systemMutex.RUnlock()

	for _, system := range sm.renderSystems {
		system.Render()
	}
}

func (sm *SystemsManager) RenderUI() {
	sm.systemMutex.RLock()
	defer sm.systemMutex.RUnlock()

	for _, system := range sm.uiRenderSystems {
		system.RenderUI()
	}
}
