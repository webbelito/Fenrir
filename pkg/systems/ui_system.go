package systems

import (
	"github.com/gen2brain/raylib-go/raygui"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/events"
)

type UISystem struct {
	ecsManager          *ecs.ECSManager
	entitiesManager     *ecs.EntitiesManager
	uiComponentsManager *ecs.UIComponentsManager
	eventsManager       *events.EventsManager
	priority            int
}

func NewUISystem(ecsM *ecs.ECSManager, p int) *UISystem {
	return &UISystem{
		ecsManager:          ecsM,
		entitiesManager:     ecsM.GetEntitiesManager(),
		uiComponentsManager: ecsM.GetUIComponentsManager(),
		eventsManager:       ecsM.GetEventsManager(),
		priority:            p,
	}
}

func (us *UISystem) RenderUI() {

	if us.ecsManager == nil || us.entitiesManager == nil || us.uiComponentsManager == nil {
		return
	}

	us.RenderUIComponents()
}

func (us *UISystem) RenderUIComponents() {

	// Render UI Panels
	us.RenderUIPanels()

	// Render UI Buttons
	us.RenderUIButtons()

	// Render UI Labels
	us.RenderUILabels()

}

func (us *UISystem) RenderUIPanels() {

	panels := us.uiComponentsManager.GetEntitiesWithComponents([]ecs.UIComponentType{ecs.UIPanelComponent})

	for _, eID := range panels {
		panelComp, panelCompExists := us.uiComponentsManager.GetComponent(eID, ecs.UIPanelComponent)

		if !panelCompExists {
			continue
		}

		panel := panelComp.(*components.UIPanel)

		if !panel.IsVisible {
			continue
		}

		raygui.Panel(panel.Bounds, panel.Title)

	}

}

func (us *UISystem) RenderUIButtons() {

	buttons := us.uiComponentsManager.GetEntitiesWithComponents([]ecs.UIComponentType{ecs.UIButtonComponent})

	for _, eID := range buttons {
		buttonComp, buttonCompExists := us.uiComponentsManager.GetComponent(eID, ecs.UIButtonComponent)

		if !buttonCompExists {
			continue
		}

		button := buttonComp.(*components.UIButton)

		if !button.IsVisible {
			continue
		}

		if raygui.Button(button.Bounds, button.Text) {
			button.OnClick(us.eventsManager)
		}

	}

}

func (us *UISystem) RenderUILabels() {

	labels := us.uiComponentsManager.GetEntitiesWithComponents([]ecs.UIComponentType{ecs.UILabelComponent})

	for _, eID := range labels {
		labelComp, labelCompExists := us.uiComponentsManager.GetComponent(eID, ecs.UILabelComponent)

		if !labelCompExists {
			continue
		}

		label := labelComp.(*components.UILabel)

		if !label.IsVisible {
			continue
		}

		raygui.Label(label.Bounds, label.Label)

	}

}

func (us *UISystem) GetPriority() int {
	return us.priority
}
