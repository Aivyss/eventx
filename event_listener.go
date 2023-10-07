package eventx

type EventListener[E any] interface {
	Trigger(entity E) error
}

func BuildEventListener[E any](trigger func(entity E) error) EventListener[E] {
	return &defaultEventListener[E]{
		InnerTrigger: trigger,
	}
}

type defaultEventListener[E any] struct {
	InnerTrigger func(entity E) error
}

func (l *defaultEventListener[E]) Trigger(entity E) error {
	return l.InnerTrigger(entity)
}
