package scenes

// SceneData is a struct that represents a scene
type SceneData struct {
	SceneName   string          `json:"scene_name"`
	Entities    []EntityData    `json:"entities"`
	Environment EnvironmentData `json:"environment"`
}

// EntityData is a struct that represents an entity
type EntityData struct {
	ID         uint64                 `json:"id"`
	Type       string                 `json:"type"`
	Position   PositionData           `json:"position"`
	Components map[string]interface{} `json:"components"`
}

// PositionData is a struct that represents a position
type PositionData struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// EnvironmentData is a struct that represents the environment
type EnvironmentData struct {
	BackgroundColor string `json:"background_color"`
	Music           string `json:"music"`
}
