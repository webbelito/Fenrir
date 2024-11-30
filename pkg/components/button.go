package components

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"github.com/webbelito/Fenrir/pkg/events"
)

type UIButton struct {
	Text      string
	Bounds    raylib.Rectangle
	IsVisible bool
}

func (b *UIButton) OnClick(em *events.EventsManager) {
	em.Dispatch("button_clicked", events.ButtonClickEvent{ButtonText: b.Text})
}
