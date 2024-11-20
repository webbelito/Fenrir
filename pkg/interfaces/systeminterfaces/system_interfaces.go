package systeminterfaces

type Updatable interface {
	Update(dt float64)
	GetPriority() int
}

type Renderable interface {
	Render()
	GetPriority() int
}
