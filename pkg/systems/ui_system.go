package systems

import (
	"github.com/gen2brain/raylib-go/raygui"
	"github.com/webbelito/Fenrir/pkg/components"
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/events"
	"github.com/webbelito/Fenrir/pkg/utils"
)

// UISystem is a system that renders UI components
type UISystem struct {
	manager  *ecs.Manager
	priority int
}

// NewUISystem creates a new UISystem
func NewUISystem(m *ecs.Manager, p int) *UISystem {
	return &UISystem{
		manager:  m,
		priority: p,
	}
}

func (us *UISystem) RenderUI() {

	if us.manager == nil {
		utils.ErrorLogger.Println("ECS Manager is nil")
		return
	}

	us.RenderUIComponents()
}

func (us *UISystem) RenderUIComponents() {

	// Render UI Panels
	us.renderUIPanels()

	// Render UI Buttons
	us.renderUIButtons()

	// Render UI Labels
	us.renderUILabels()

}

/*
renderUIPanels renders all UIPanel components
*/
func (us *UISystem) renderUIPanels() {

	// Get all entities with UIPanelComponent
	panels := us.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.UIPanelComponent})

	// Iterate over all entities with UIPanelComponent
	for _, eID := range panels {

		// Get the UIPanelComponent
		panelComp, panelCompExists := us.manager.GetComponent(eID, ecs.UIPanelComponent)

		// If the UIPanelComponent does not exist, skip to the next entity
		if !panelCompExists {
			utils.ErrorLogger.Println("UIPanelComponent does not exist")
			continue
		}

		// Type assert the UIPanelComponent
		panel, ok := panelComp.(*components.UIPanel)

		// If the type assertion fails, skip to the next entity
		if !ok {
			utils.ErrorLogger.Println("UIPanel type assertion failed")
			continue
		}

		if !panel.IsVisible {
			continue
		}

		// Render the panel
		raygui.Panel(panel.Bounds, panel.Title)

	}

}

/*
renderUIButtons renders all UIButton components
*/
func (us *UISystem) renderUIButtons() {

	// Get all entities with UIButtonComponent
	buttons := us.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.UIButtonComponent})

	// Iterate over all entities with UIButtonComponent
	for _, eID := range buttons {

		// Get the UIButtonComponent
		buttonComp, buttonCompExists := us.manager.GetComponent(eID, ecs.UIButtonComponent)

		// If the UIButtonComponent does not exist, skip to the next entity
		if !buttonCompExists {
			utils.ErrorLogger.Println("UIButtonComponent does not exist")
			continue
		}

		// Type assert the UIButtonComponent
		button, ok := buttonComp.(*components.UIButton)

		// If the type assertion fails, skip to the next entity
		if !ok {
			utils.ErrorLogger.Println("UIButton type assertion failed")
			continue
		}

		// If the button is not visible, skip to the next entity
		if !button.IsVisible {
			continue
		}

		// Render the button
		if raygui.Button(button.Bounds, button.Text) {

			// If the button is clicked, call the OnClick function
			us.manager.DispatchEvent("button_clicked", events.ButtonClickEvent{ButtonText: button.Text})
		}
	}
}

/*
renderUILabels renders all UILabel components
*/
func (us *UISystem) renderUILabels() {

	// Get all entities with UILabelComponent
	labels := us.manager.GetEntitiesWithComponents([]ecs.ComponentType{ecs.UILabelComponent})

	// Iterate over all entities with UILabelComponent
	for _, eID := range labels {

		// Get the UILabelComponent
		labelComp, labelCompExists := us.manager.GetComponent(eID, ecs.UILabelComponent)

		if !labelCompExists {
			utils.ErrorLogger.Println("UILabelComponent does not exist")
			continue
		}

		// Type assert the UILabelComponent
		label, ok := labelComp.(*components.UILabel)

		// If the type assertion fails, skip to the next entity
		if !ok {
			utils.ErrorLogger.Println("UILabel type assertion failed")
			continue
		}

		// If the label is not visible, skip to the next entity
		if !label.IsVisible {
			continue
		}

		// Render the label
		raygui.Label(label.Bounds, label.Label)

	}
}

/*
GetPriority returns the priority of the system
*/
func (us *UISystem) GetPriority() int {
	return us.priority
}
