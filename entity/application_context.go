package entity

import (
	"context"
	"fmt"
	"github.com/aivyss/eventx/common"
	"reflect"
	"sync"
)

const (
	DefaultEventChannelBufferSize = 5
	DefaultEventProcessPoolSize   = 10
	DefaultMultiEventMode         = true
)

type ApplicationContext struct {
	once                   sync.Once
	innerContext           context.Context
	innerContextCancel     context.CancelFunc
	multiEventMode         bool
	eventChannel           chan EventRunner
	eventChannelBufferSize int
	eventProcessPoolSize   int
	listenerdMap           common.MultiMap[reflect.Type, any]
}

func NewApplicationContext(eventChannelBufferSize int, eventProcessPoolSize int, multiEventMode bool) *ApplicationContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &ApplicationContext{
		innerContext:           ctx,
		innerContextCancel:     cancel,
		multiEventMode:         multiEventMode,
		eventChannel:           make(chan EventRunner, eventChannelBufferSize),
		eventChannelBufferSize: eventChannelBufferSize,
		eventProcessPoolSize:   eventProcessPoolSize,
		listenerdMap:           common.NewMultiMap[reflect.Type, any](),
	}
}

func (ctx *ApplicationContext) QueueEventRunner(runner EventRunner) {
	ctx.eventChannel <- runner
}

func (ctx *ApplicationContext) ConsumeEventRunner() {
	ctx.once.Do(func() {
		setting := "default"
		if ctx.eventChannelBufferSize != DefaultEventChannelBufferSize ||
			ctx.eventProcessPoolSize != DefaultEventProcessPoolSize ||
			ctx.multiEventMode {
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
					case runner := <-ctx.eventChannel:
						runner()
					}
				}
			}(ctx.innerContext)
		}
	})
}

func (ctx *ApplicationContext) GetEventListener(typeVal reflect.Type) []any {
	listeners := ctx.listenerdMap.Get(typeVal)
	return listeners
}

func (ctx *ApplicationContext) RegisterEventListener(typeVal reflect.Type, eventListener any) {
	ctx.listenerdMap.Put(typeVal, eventListener)
}

func (ctx *ApplicationContext) Close() {
	ctx.innerContextCancel()
}

func (ctx *ApplicationContext) IsMultiMode() bool {
	return ctx.multiEventMode
}
