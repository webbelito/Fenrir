package resources

import (
	"fmt"
	"sync"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// ResourcesManager is a struct that holds all the resources that the game uses.
type ResourcesManager struct {
	textures map[string]raylib.Texture2D
	mutex    sync.RWMutex
}

// NewResourceManager creates a new ResourceManager.
func NewResourceManager() *ResourcesManager {
	return &ResourcesManager{
		textures: make(map[string]raylib.Texture2D),
	}
}

// LoadTexture loads a texture from a file and stores it in the ResourceManager.
func (rm *ResourcesManager) LoadTexture(path string) (raylib.Texture2D, error) {
	rm.mutex.Lock()

	// Check if the texture is already loaded
	if texture, texExists := rm.textures[path]; texExists {
		rm.mutex.Unlock()
		return texture, nil
	}

	rm.mutex.Unlock()

	// Load the texture
	texture := raylib.LoadTexture(path)
	if texture.ID == 0 {
		return texture, fmt.Errorf("failed to load texture: %s", path)
	}

	// Store the loaded texture
	rm.mutex.Lock()
	rm.textures[path] = texture
	rm.mutex.Unlock()

	return texture, nil
}

func (rm *ResourcesManager) GetTexture(path string) (raylib.Texture2D, bool) {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	texture, texExists := rm.textures[path]
	return texture, texExists
}

func (rm *ResourcesManager) UnloadTextures() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	// Free all textures
	for _, texture := range rm.textures {
		raylib.UnloadTexture(texture)
	}

	// Clear the map
	rm.textures = make(map[string]raylib.Texture2D)
}
