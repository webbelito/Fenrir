package systeminterfaces

type Updatable interface {
	Update(dt float64)
}

type Renderable interface {
	Render()
}
