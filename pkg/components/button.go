package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/events"
)

// UIButton is a component that holds a UI button
type UIButton struct {
	Text      string
	Bounds    raylib.Rectangle
	IsVisible bool
}

// OnClick is a method that dispatches a button click event
func (b *UIButton) OnClick(em *events.EventsManager) {
	em.Dispatch("button_clicked", events.ButtonClickEvent{ButtonText: b.Text})
}
