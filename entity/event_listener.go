package entity

type EventListener[E any] interface {
	Trigger(entity E) error
}

type CatchErrEventListener[E any] interface {
	EventListener[E]
	Catch(err error)
}

type SuccessEventListener[E any] interface {
	EventListener[E]
	Then(entity E)
}

type CallbackEventListener[E any] interface {
	EventListener[E]
	SuccessEventListener[E]
	CatchErrEventListener[E]
}

type TriggerFunc[E any] func(entity E) error
type ThenFunc[E any] func(entity E)
type CatchFunc func(err error)
