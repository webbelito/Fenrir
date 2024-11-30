package editor

import (
	"github.com/webbelito/Fenrir/pkg/ecs"
)

// EditorManager is a manager for the editor
type EditorManager struct {
	editor    *Editor
	isVisible bool
	priority  int
}

// NewEditorManager creates a new EditorManager
func NewEditorManager(m *ecs.Manager, p int) *EditorManager {

	// Create a new editor
	e := NewEditor(m)

	return &EditorManager{
		editor:    e,
		isVisible: false,
		priority:  p,
	}
}

// ToggleVisibility toggles the visibility of the editor
func (em *EditorManager) ToggleVisibility() {
	em.isVisible = !em.isVisible
}

func (em *EditorManager) Update(dt float64) {

	// If the editor is not visible, return
	if !em.isVisible {
		return
	}

	// Update the editor
	em.editor.Update(dt)
}

func (em *EditorManager) Render() {

	// If the editor is not visible, return
	if !em.isVisible {
		return
	}

	// Render the editor
	em.editor.Render()
}

// GetPriority returns the priority of the editor
func (em *EditorManager) GetPriority() int {
	return em.priority
}
