package context

import "github.com/aivyss/eventx/entity"

type EventListenerDispenseChannel struct {
	DispensePoolSize   int
	DispenseBufferSize int
	DispenseChannel    chan entity.EventSet
}
