package entity

func BuildEventListenerWithCallback[E any](
	trigger func(entity E) error,
	then func(entity E),
	catch func(err error),
) EventListener[E] {
	return &callbackEventListener[E]{
		InnerTrigger: trigger,
		InnerThen:    then,
		InnerCatch:   catch,
	}
}

type callbackEventListener[E any] struct {
	InnerTrigger TriggerFunc[E]
	InnerThen    ThenFunc[E]
	InnerCatch   CatchFunc
}

func (l *callbackEventListener[E]) Then(entity E) {
	l.InnerThen(entity)
}

func (l *callbackEventListener[E]) Catch(err error) {
	l.InnerCatch(err)
}

func (l *callbackEventListener[E]) Trigger(entity E) error {
	return l.InnerTrigger(entity)
}
