package system_interfaces

type Updatable interface {
	Update(dt float64)
}

type Renderable interface {
	Render()
}
