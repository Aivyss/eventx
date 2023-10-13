package entity

func BuildEventListener[E any](trigger func(entity E) error) EventListener[E] {
	return &defaultEventListener[E]{
		InnerTrigger: trigger,
	}
}

type defaultEventListener[E any] struct {
	InnerTrigger TriggerFunc[E]
}

func (l *defaultEventListener[E]) Then(_ E) {}

func (l *defaultEventListener[E]) Catch(_ error) {}

func (l *defaultEventListener[E]) Trigger(entity E) error {
	return l.InnerTrigger(entity)
}
