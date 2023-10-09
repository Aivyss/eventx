package context

import "github.com/aivyss/eventx/entity"

type EventChannel struct {
	ChannelBufferSize int
	ProcessPoolSize   int
	Channel           chan entity.EventRunner
}
