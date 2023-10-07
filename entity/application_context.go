package entity

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

const DefaultEventChannelBufferSize = 5
const DefaultEventProcessPoolSize = 10

type ApplicationContext struct {
	once                   sync.Once
	innerContext           context.Context
	innerContextCancel     context.CancelFunc
	eventChannel           chan EventRunner
	eventChannelBufferSize int
	eventProcessPoolSize   int
	listenerMap            map[reflect.Type]any
}

func NewApplicationContext(eventChannelBufferSize int, eventProcessPoolSize int) *ApplicationContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &ApplicationContext{
		innerContext:           ctx,
		innerContextCancel:     cancel,
		eventChannel:           make(chan EventRunner, eventChannelBufferSize),
		eventChannelBufferSize: eventChannelBufferSize,
		eventProcessPoolSize:   eventProcessPoolSize,
		listenerMap:            map[reflect.Type]any{},
	}
}

func (ctx *ApplicationContext) QueueEventRunner(runner EventRunner) {
	ctx.eventChannel <- runner
}

func (ctx *ApplicationContext) ConsumeEventRunner() {
	ctx.once.Do(func() {
		setting := "default"
		if ctx.eventChannelBufferSize != DefaultEventChannelBufferSize ||
			ctx.eventProcessPoolSize != DefaultEventProcessPoolSize {
			setting = "customized"
		}
		fmt.Println(fmt.Sprintf(
			"[eventx] eventx event channel is running (setting: %s)",
			setting,
		))
		fmt.Println(fmt.Sprintf(
			"[eventx] EventChannelSize: %d, EventProcessPoolSize: %d",
			ctx.eventChannelBufferSize,
			ctx.eventProcessPoolSize,
		))

		for i := 0; i < ctx.eventProcessPoolSize; i++ {
			go func(innerContext context.Context) {
				for {
					select {
					case <-innerContext.Done():
						fmt.Println("[eventx] End of event process pool...")
						return
					default:
						(<-ctx.eventChannel)()
					}
				}
			}(ctx.innerContext)
		}
	})
}

func (ctx *ApplicationContext) GetEventListener(typeVal reflect.Type) (any, bool) {
	unspecifiedEventListener, ok := ctx.listenerMap[typeVal]

	return unspecifiedEventListener, ok
}

func (ctx *ApplicationContext) RegisterEventListener(typeVal reflect.Type, eventListener any) {
	ctx.listenerMap[typeVal] = eventListener
}

func (ctx *ApplicationContext) Close() {
	ctx.innerContextCancel()
}
