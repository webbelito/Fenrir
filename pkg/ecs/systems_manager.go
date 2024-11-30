package ecs

import (
	"sync"

	systeminterfaces "github.com/webbelito/Fenrir/pkg/interfaces/systeminterfaces"
	"github.com/webbelito/Fenrir/pkg/utils"
)

// SystemsManager is a struct that manages all systems
type SystemsManager struct {
	logicSystems    []systeminterfaces.UpdatableSystemInterface
	renderSystems   []systeminterfaces.RenderableSystemInterface
	uiRenderSystems []systeminterfaces.UIRenderableSystemInterface
	systemMutex     sync.RWMutex
}

// NewSystemsManager creates a new SystemsManager
func NewSystemsManager() *SystemsManager {
	return &SystemsManager{
		logicSystems:    []systeminterfaces.UpdatableSystemInterface{},
		renderSystems:   []systeminterfaces.RenderableSystemInterface{},
		uiRenderSystems: []systeminterfaces.UIRenderableSystemInterface{},
	}
}

// AddLogicSystem adds a logic system to the SystemsManager
func (sm *SystemsManager) AddLogicSystem(system systeminterfaces.UpdatableSystemInterface, priority int) {

	// Check if the system is already in the list
	inserted := false

	// TODO: Implement priority

	// Insert the system in the correct position
	for i, existingSystem := range sm.logicSystems {

		// If the existing system has a higher priority than the new system, insert the new system before the existing system
		if existingSystem.GetPriority() > system.GetPriority() {
			sm.logicSystems = append(sm.logicSystems[:i], append([]systeminterfaces.UpdatableSystemInterface{system}, sm.logicSystems[i:]...)...)

			inserted = true

			utils.InfoLogger.Printf("Added logic system: %T\n", system)

			break
		}
	}

	// If the system was not inserted, append it to the end of the list
	if !inserted {
		sm.logicSystems = append(sm.logicSystems, system)
		utils.InfoLogger.Printf("Added logic system: %T\n", system)
	}

}

// RemoveLogicSystem removes a logic system from the SystemsManager
func (sm *SystemsManager) RemoveLogicSystem(system systeminterfaces.UpdatableSystemInterface) {

	sm.systemMutex.Lock()
	defer sm.systemMutex.Unlock()

	// Find the system in the list and remove it
	for i, sys := range sm.logicSystems {

		// If the system is found, remove it from the list
		if sys == system {
			sm.logicSystems = append(sm.logicSystems[:i], sm.logicSystems[i+1:]...)
			utils.InfoLogger.Printf("Removed logic system: %T\n", system)
			break
		}
	}
}

// AddRenderSystem adds a render system to the SystemsManager
func (sm *SystemsManager) AddRenderSystem(system systeminterfaces.RenderableSystemInterface, priority int) {

	//sm.systemMutex.Lock()
	//defer sm.systemMutex.Unlock()

	// Check if the system is already in the list
	inserted := false

	// TODO: Implement priority

	// Insert the system in the correct position
	for i, existingSystem := range sm.renderSystems {

		// If the existing system has a higher priority than the new system, insert the new system before the existing system
		if existingSystem.GetPriority() > system.GetPriority() {
			sm.renderSystems = append(sm.renderSystems[:i], append([]systeminterfaces.RenderableSystemInterface{system}, sm.renderSystems[i:]...)...)
			inserted = true
			utils.InfoLogger.Printf("Added render system: %T\n", system)
			break
		}
	}

	// If the system was not inserted, append it to the end of the list
	if !inserted {
		sm.renderSystems = append(sm.renderSystems, system)
		utils.InfoLogger.Printf("Added render system: %T\n", system)
	}
}

// RemoveRenderSystem removes a render system from the SystemsManager
func (sm *SystemsManager) RemoveRenderSystem(system systeminterfaces.RenderableSystemInterface) {

	sm.systemMutex.Lock()
	defer sm.systemMutex.Unlock()

	// Find the system in the list and remove it
	for i, sys := range sm.renderSystems {

		// If the system is found, remove it from the list
		if sys == system {
			sm.renderSystems = append(sm.renderSystems[:i], sm.renderSystems[i+1:]...)
			utils.InfoLogger.Printf("Removed render system: %T\n", system)
			break
		}
	}
}

// AddUIRenderSystem adds a UI render system to the SystemsManager
func (sm *SystemsManager) AddUIRenderSystem(system systeminterfaces.UIRenderableSystemInterface, priority int) {

	sm.systemMutex.Lock()
	defer sm.systemMutex.Unlock()

	// Check if the system is already in the list
	inserted := false

	// TODO: Implement priority

	// Insert the system in the correct position
	for i, existingSystem := range sm.uiRenderSystems {

		// If the existing system has a higher priority than the new system, insert the new system before the existing system
		if existingSystem.GetPriority() > system.GetPriority() {
			sm.uiRenderSystems = append(sm.uiRenderSystems[:i], append([]systeminterfaces.UIRenderableSystemInterface{system}, sm.uiRenderSystems[i:]...)...)
			inserted = true
			utils.InfoLogger.Printf("Added UI render system: %T\n", system)
			break
		}
	}

	// If the system was not inserted, append it to the end of the list
	if !inserted {
		sm.uiRenderSystems = append(sm.uiRenderSystems, system)
	}
}

// RemoveUIRenderSystem removes a UI render system from the SystemsManager
func (sm *SystemsManager) RemoveUIRenderSystem(system systeminterfaces.UIRenderableSystemInterface) {

	sm.systemMutex.Lock()
	defer sm.systemMutex.Unlock()

	// Find the system in the list and remove it
	for i, sys := range sm.uiRenderSystems {

		// If the system is found, remove it from the list
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

		// Check if the system is a CameraSystem
		if cameraSys, ok := sys.(systeminterfaces.CameraSystemInterface); ok {
			return cameraSys, true
		}
	}

	// Check Render Systems
	for _, sys := range sm.renderSystems {

		// Check if the system is a CameraSystem
		if cameraSys, ok := sys.(systeminterfaces.CameraSystemInterface); ok {
			return cameraSys, true
		}
	}

	return nil, false
}

func (sm *SystemsManager) Update(dt float64) {

	sm.systemMutex.RLock()
	defer sm.systemMutex.RUnlock()

	// Update all logic systems
	for _, system := range sm.logicSystems {
		system.Update(dt)
	}
}

func (sm *SystemsManager) Render() {

	sm.systemMutex.RLock()
	defer sm.systemMutex.RUnlock()

	// Render all render systems
	for _, system := range sm.renderSystems {
		system.Render()
	}
}

func (sm *SystemsManager) RenderUI() {

	sm.systemMutex.RLock()
	defer sm.systemMutex.RUnlock()

	// Render all UI render systems
	for _, system := range sm.uiRenderSystems {
		system.RenderUI()
	}
}
