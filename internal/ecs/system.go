package ecs

type System interface {
	Update(deltaTime float32, manager *Manager)
}
