package entity

type EventListener[E any] interface {
	Trigger(entity E) error
}

type EventFunc[E any] func(entity E) error

func BuildEventListener[E any](trigger func(entity E) error) EventListener[E] {
	return &defaultEventListener[E]{
		InnerTrigger: trigger,
	}
}

type defaultEventListener[E any] struct {
	InnerTrigger EventFunc[E]
}

func (l *defaultEventListener[E]) Trigger(entity E) error {
	return l.InnerTrigger(entity)
}
