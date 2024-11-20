package editor

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
)

type EditorManager struct {
	editor    *Editor
	isVisible bool
	priority  int
}

func NewEditorManager(ecsM *ecs.ECSManager, p int) *EditorManager {

	e := NewEditor(ecsM)

	return &EditorManager{
		editor:    e,
		isVisible: false,
		priority:  p,
	}
}

func (em *EditorManager) ToggleVisibility() {
	em.isVisible = !em.isVisible
}

func (em *EditorManager) Update(dt float64) {
	if !em.isVisible {
		return
	}

	em.editor.Update(dt)
}

func (em *EditorManager) Render() {
	if !em.isVisible {
		return
	}
	em.editor.Render()
}

func (em *EditorManager) GetPriority() int {
	return em.priority
}
