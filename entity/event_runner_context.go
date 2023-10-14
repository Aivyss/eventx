package entity

import "sync"

type EventContext interface {
	IsRunnable() bool
	Cancel() bool
	IsDone() bool
}

type EventRunnerContextImpl struct {
	sync.Mutex
	Runnable bool
	Done     bool
}

func NewEventRunnerContext() *EventRunnerContextImpl {
	return &EventRunnerContextImpl{Runnable: true, Done: false}
}

func (c *EventRunnerContextImpl) IsRunnable() bool {
	runnable := false

	c.Lock()
	runnable = c.Runnable && (!c.Done)
	c.Unlock()

	return runnable
}

func (c *EventRunnerContextImpl) IsRunnableInternal() bool {
	return c.Runnable && (!c.Done)
}

func (c *EventRunnerContextImpl) Cancel() bool {
	done := false

	c.Lock()
	c.Runnable = false
	done = c.Done
	c.Unlock()

	return done
}

func (c *EventRunnerContextImpl) IsDone() bool {
	finished := false

	c.Lock()
	finished = c.Done
	c.Unlock()

	return finished
}
