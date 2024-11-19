package scenes

type SceneData struct {
	SceneName   string          `json:"scene_name"`
	Entities    []EntityData    `json:"entities"`
	Environment EnvironmentData `json:"environment"`
}

type EntityData struct {
	ID         uint64                 `json:"id"`
	Type       string                 `json:"type"`
	Position   PositionData           `json:"position"`
	Components map[string]interface{} `json:"components"`
}

type PositionData struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type EnvironmentData struct {
	BackgroundColor string `json:"background_color"`
	Music           string `json:"music"`
}
