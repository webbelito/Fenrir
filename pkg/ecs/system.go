package ecs

type System interface {
	Update(dt float64, em *EntitiesManager, cm *ComponentsManager)
}
