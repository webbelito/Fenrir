package scenes

import (
	"encoding/json"
	"os"
)

func LoadSceneData(sceneFilePath string) (*SceneData, error) {

	data, err := os.ReadFile(sceneFilePath)
	if err != nil {
		return nil, err
	}

	var sceneData SceneData
	err = json.Unmarshal(data, &sceneData)
	if err != nil {
		return nil, err
	}

	return &sceneData, nil

}
