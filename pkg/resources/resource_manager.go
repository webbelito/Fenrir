package resources

import (
	"fmt"
	"sync"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

// ResourcesManager is a struct that holds all the resources that the game uses.
type ResourcesManager struct {
	textures map[string]raylib.Texture2D
	sounds   map[string]raylib.Sound
	mutex    sync.RWMutex
}

// NewResourcesManager creates a new ResourceManager.
func NewResourcesManager() *ResourcesManager {
	return &ResourcesManager{
		textures: make(map[string]raylib.Texture2D),
		sounds:   make(map[string]raylib.Sound),
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

// LoadSound loads a sound from a file and stores it in the ResourceManager.
func (rm *ResourcesManager) LoadSound(path string) (raylib.Sound, error) {
	rm.mutex.RLock()
	sound, soundExists := rm.sounds[path]
	rm.mutex.RUnlock()

	if soundExists {
		return sound, nil
	}

	loadedSound := raylib.LoadSound(path)
	if loadedSound.FrameCount == 0 {
		return loadedSound, fmt.Errorf("failed to load: %v", path)
	}

	rm.mutex.Lock()
	rm.sounds[path] = loadedSound
	rm.mutex.Unlock()

	return loadedSound, nil
}

func (rm *ResourcesManager) UnloadAll() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	for _, texture := range rm.textures {
		raylib.UnloadTexture(texture)
	}
	rm.textures = make(map[string]raylib.Texture2D)

	for _, sound := range rm.sounds {
		raylib.UnloadSound(sound)
	}
	rm.sounds = make(map[string]raylib.Sound)
}
