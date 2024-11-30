package systems

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
	"github.com/webbelito/Fenrir/pkg/events"
)

type EventsListenerSystem struct {
	ecsManager    *ecs.ECSManager
	eventsManager *events.EventsManager
	priority      int
}

func NewEventsListenerSystem(ecsM *ecs.ECSManager, p int) *EventsListenerSystem {
	els := &EventsListenerSystem{
		ecsManager:    ecsM,
		eventsManager: ecsM.GetEventsManager(),
		priority:      p,
	}

	els.eventsManager.Subscribe("button_clicked", els.OnButtonClick)

	return els

}

func (els *EventsListenerSystem) OnButtonClick(event events.Event) {

	switch e := event.(type) {
	case events.ButtonClickEvent:
		switch e.ButtonText {
		case "Start Game":
			els.ecsManager.GetEventsManager().Dispatch("change_scene", events.SceneChangeEvent{ScenePath: "assets/scenes/game_scene.json"})
		case "Options":
			els.ecsManager.GetEventsManager().Dispatch("change_scene", events.SceneChangeEvent{ScenePath: "assets/scenes/options_scene.json"})
		case "Exit":
			els.ecsManager.GetEventsManager().Dispatch("exit_game", events.ExitGameEvent{ShouldExitGame: true})
		}
	}
}

func (els *EventsListenerSystem) GetPriority() int {
	return els.priority
}

func (els *EventsListenerSystem) Update(dt float64) {
	// Do nothing
}
