package entity

type EventSet interface {
	Runner() func()
	Context() *EventRunnerContextImpl
}

type EventSetImpl[E any] struct {
	EventListener EventListener[E]
	Entity        E
	Ctx           *EventRunnerContextImpl
}

func NewEventSet[E any](listener EventListener[E], entity E) EventSet {
	return &EventSetImpl[E]{
		EventListener: listener,
		Entity:        entity,
		Ctx:           NewEventRunnerContext(),
	}
}

func (s *EventSetImpl[E]) Runner() func() {
	err := s.EventListener.Trigger(s.Entity)
	if err != nil {
		el, ok := s.EventListener.(CatchErrEventListener[E])

		if ok {
			return func() {
				el.Catch(err)
			}
		}

		return nil
	}

	el, ok := s.EventListener.(SuccessEventListener[E])
	if ok {
		return func() {
			el.Then(s.Entity)
		}
	}

	return nil
}

func (s *EventSetImpl[E]) Context() *EventRunnerContextImpl {
	return s.Ctx
}
