package events

// Event is an interface that all events must implement
type Event interface{}

// ButtonClickEvent represents a button click event
type ButtonClickEvent struct {
	ButtonText string
}

// SceneChangeEvent represents a scene change event
type SceneChangeEvent struct {
	ScenePath string
}

// ExitGameEvent represents an exit game event
type ExitGameEvent struct {
	ShouldExitGame bool
}
