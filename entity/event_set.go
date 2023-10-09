package entity

type EventSet interface {
	Runner()
}

type EventSetImpl[E any] struct {
	EventListener EventListener[E]
	Entity        E
}

func NewEventSet[E any](listener EventListener[E], entity E) EventSet {
	return &EventSetImpl[E]{
		EventListener: listener,
		Entity:        entity,
	}
}

func (s *EventSetImpl[E]) Runner() {
	_ = s.EventListener.Trigger(s.Entity)
}
