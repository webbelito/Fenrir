package systems

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/events"
)

// EventsListenerSystem is a system that listens for events
type EventsListenerSystem struct {
	manager  *ecs.Manager
	priority int
}

// NewEventsListenerSystem creates a new EventsListenerSystem
func NewEventsListenerSystem(m *ecs.Manager, p int) *EventsListenerSystem {
	els := &EventsListenerSystem{
		manager:  m,
		priority: p,
	}

	// Subscribe to button clicked events
	els.manager.SubscribeEvent("button_clicked", els.OnButtonClick)

	return els
}

// OnButtonClick handles button click events
func (els *EventsListenerSystem) OnButtonClick(event events.Event) {

	// Check the type of the event
	switch e := event.(type) {

	// If the event is a ButtonClickEvent
	case events.ButtonClickEvent:

		// Check the button text, and dispatch the appropriate event
		switch e.ButtonText {
		case "Start Game":
			els.manager.DispatchEvent("change_scene", events.SceneChangeEvent{ScenePath: "assets/scenes/game_scene.json"})
		case "Options":
			els.manager.DispatchEvent("change_scene", events.SceneChangeEvent{ScenePath: "assets/scenes/options_scene.json"})
		case "Exit":
			els.manager.DispatchEvent("exit_game", events.ExitGameEvent{ShouldExitGame: true})
		}
	}
}

/*
GetPriority returns the priority of the system
*/
func (els *EventsListenerSystem) GetPriority() int {
	return els.priority
}

func (els *EventsListenerSystem) Update(dt float64) {
	// Do nothing
}
