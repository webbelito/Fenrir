package ecs

import (
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
