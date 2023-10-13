package entity

type EventListener[E any] interface {
	Trigger(entity E) error
	Then(entity E)
	Catch(err error)
}

type TriggerFunc[E any] func(entity E) error
type ThenFunc[E any] func(entity E)
type CatchFunc func(err error)
