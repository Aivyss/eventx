package entity

func BuildSuccessEventListener[E any](
	trigger func(entity E) error,
	then func(entity E),
) EventListener[E] {
	return &successEventListener[E]{
		InnerTrigger: trigger,
		InnerThen:    then,
	}
}

type successEventListener[E any] struct {
	InnerTrigger TriggerFunc[E]
	InnerThen    ThenFunc[E]
}

func (l *successEventListener[E]) Then(entity E) {
	l.InnerThen(entity)
}

func (l *successEventListener[E]) Trigger(entity E) error {
	return l.InnerTrigger(entity)
}
