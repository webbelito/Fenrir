package scenes

import (
	"encoding/json"
	"os"

	"github.com/webbelito/Fenrir/pkg/utils"
)

// LoadSceneData loads a scene from a file
func LoadSceneData(sceneFilePath string) (*SceneData, error) {

	// Read the scene file
	data, err := os.ReadFile(sceneFilePath)
	if err != nil {
		utils.ErrorLogger.Println("Failed to read scene file: ", err)
		return nil, err
	}

	// Declare a SceneData struct
	var sceneData SceneData

	// Unmarshal the scene data
	err = json.Unmarshal(data, &sceneData)
	if err != nil {
		utils.ErrorLogger.Println("Failed to unmarshal scene data: ", err)
		return nil, err
	}

	return &sceneData, nil

}
