package context

import (
	"context"
	"fmt"
	"github.com/aivyss/eventx/common"
	"github.com/aivyss/eventx/entity"
	"reflect"
	"sync"
)

const (
	DefaultEventChannelBufferSize = 5
	DefaultEventProcessPoolSize   = 10
	DefaultMultiEventMode         = true
)

type ApplicationContext struct {
	once                sync.Once
	innerContext        context.Context
	innerContextCancel  context.CancelFunc
	eventChannel        *EventChannel
	eventListenerConfig *EventListenerConfig
	*EventListenerDispenseChannel
}

func NewApplicationContext(eventChannelBufferSize int, eventProcessPoolSize int, multiEventMode bool) *ApplicationContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &ApplicationContext{
		innerContext:       ctx,
		innerContextCancel: cancel,
		eventChannel: &EventChannel{
			Channel:           make(chan entity.EventRunner, eventChannelBufferSize),
			ChannelBufferSize: eventChannelBufferSize,
			ProcessPoolSize:   eventProcessPoolSize,
		},
		eventListenerConfig: &EventListenerConfig{
			MultiEventMode: multiEventMode,
			ListenerMap:    common.NewMultiMap[reflect.Type, any](),
		},
		EventListenerDispenseChannel: &EventListenerDispenseChannel{
			DispenseBufferSize: 1,
			DispensePoolSize:   3,
			DispenseChannel:    make(chan entity.EventSet, 1),
		},
	}
}

func (ctx *ApplicationContext) QueueEventRunner(runner entity.EventRunner) {
	ctx.eventChannel.Channel <- runner
}

func (ctx *ApplicationContext) QueueEventSet(set entity.EventSet) {
	ctx.DispenseChannel <- set
}

func (ctx *ApplicationContext) ConsumeEventRunner() {
	ctx.once.Do(func() {
		setting := "default"
		if ctx.eventChannel.ChannelBufferSize != DefaultEventChannelBufferSize ||
			ctx.eventChannel.ProcessPoolSize != DefaultEventProcessPoolSize ||
			ctx.eventListenerConfig.MultiEventMode {
			setting = "customized"
		}
		fmt.Println(fmt.Sprintf(
			"[eventx] eventx event channel is running (setting: %s)",
			setting,
		))
		fmt.Println(fmt.Sprintf(
			"[eventx] EventChannelSize: %d, EventProcessPoolSize: %d",
			ctx.eventChannel.ChannelBufferSize,
			ctx.eventChannel.ProcessPoolSize,
		))

		for i := 0; i < ctx.DispensePoolSize; i++ {
			go func(innerContext context.Context) {
				for {
					select {
					case <-innerContext.Done():
						fmt.Println("[eventx] End of event listener dispense pool...")
						return
					case eventSet := <-ctx.DispenseChannel:
						ctx.eventChannel.Channel <- eventSet.Runner
					}
				}
			}(ctx.innerContext)
		}

		for i := 0; i < ctx.eventChannel.ProcessPoolSize; i++ {
			go func(innerContext context.Context) {
				for {
					select {
					case <-innerContext.Done():
						fmt.Println("[eventx] End of event process pool...")
						return
					case runner := <-ctx.eventChannel.Channel:
						runner()
					}
				}
			}(ctx.innerContext)
		}
	})
}

func (ctx *ApplicationContext) GetEventListener(typeVal reflect.Type) []any {
	listeners := ctx.eventListenerConfig.ListenerMap.Get(typeVal)
	return listeners
}

func (ctx *ApplicationContext) RegisterEventListener(typeVal reflect.Type, eventListener any) {
	ctx.eventListenerConfig.ListenerMap.Put(typeVal, eventListener)
}

func (ctx *ApplicationContext) Close() {
	ctx.innerContextCancel()
}

func (ctx *ApplicationContext) IsMultiMode() bool {
	return ctx.eventListenerConfig.MultiEventMode
}
