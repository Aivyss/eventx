package entity

func BuildCatchErrEventListener[E any](
	trigger func(entity E) error,
	catch func(err error),
) EventListener[E] {
	return &catchErrEventListener[E]{
		InnerTrigger: trigger,
		InnerCatch:   catch,
	}
}

type catchErrEventListener[E any] struct {
	InnerTrigger TriggerFunc[E]
	InnerCatch   CatchFunc
}

func (l *catchErrEventListener[E]) Catch(err error) {
	l.InnerCatch(err)
}

func (l *catchErrEventListener[E]) Trigger(entity E) error {
	return l.InnerTrigger(entity)
}
