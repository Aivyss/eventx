package entity

type EventSet interface {
	Runner() func()
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

func (s *EventSetImpl[E]) Runner() func() {
	err := s.EventListener.Trigger(s.Entity)
	if err != nil {
		return func() {
			s.EventListener.Catch(err)
		}
	}

	return func() {
		s.EventListener.Then(s.Entity)
	}
}
