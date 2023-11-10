package context

import (
	"context"
	"fmt"
	"github.com/aivyss/eventx/entity"
	"github.com/aivyss/typex"
	"reflect"
	"sync"
)

const (
	DefaultEventChannelBufferSize = 5
	DefaultEventProcessPoolSize   = 10
	DefaultMultiEventMode         = true
)

// ApplicationContext
//
// The context needed to operate the entirety of `eventx` application.
//
// There should be only one variable declared globally in the application (eventx/global_variables.go/appContext).
// (For internal usage within `eventx`)
//
// Users can create ApplicationContext as well with ApplicationContext{} and NewApplicationContext function,
// but cannot assign values to its internal fields and cannot use it for `eventx` application purposes.
type ApplicationContext struct {
	// once
	//
	// Used for the sole purpose of running the ConsumeEventRunner method only once.
	once sync.Once
	// innerContext
	//
	// The context of the event pool operated within the ApplicationContext.
	//
	// eventx.Close => It ends when the Close method is executed.
	innerContext context.Context
	// innerContextCancel
	//
	// Ends the innerContext, causing the event pool to terminate.
	innerContextCancel context.CancelFunc
	// eventChannel manages the channels used for executing events.
	eventChannel *EventChannel
	// eventListenerConfig manages the state of entity.EventListener.
	eventListenerConfig *EventListenerConfig
	// EventListenerDispenseChannel is an intermediate layer for event listener processing distribution.
	eventListenerDispenseChannel *EventListenerDispenseChannel
}

// NewApplicationContext
//
// Creates an ApplicationContext.
//
// Users can create ApplicationContext as well with ApplicationContext{} and the NewApplicationContext function,
// but cannot assign values to its internal fields and cannot use it for `eventx` application purposes.
func NewApplicationContext(eventChannelBufferSize int, eventProcessPoolSize int, multiEventMode bool) *ApplicationContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &ApplicationContext{
		innerContext:       ctx,
		innerContextCancel: cancel,
		eventChannel: &EventChannel{
			Channel:           make(chan entity.EventRunner, eventChannelBufferSize),
			AfterChannel:      make(chan entity.EventAfterRunner, eventChannelBufferSize),
			ChannelBufferSize: eventChannelBufferSize,
			ProcessPoolSize:   eventProcessPoolSize,
		},
		eventListenerConfig: &EventListenerConfig{
			MultiEventMode: multiEventMode,
			ListenerMap:    typex.NewMultiMap[reflect.Type, any](),
		},
		eventListenerDispenseChannel: &EventListenerDispenseChannel{
			DispenseBufferSize: 1,
			DispensePoolSize:   3,
			DispenseChannel:    make(chan entity.EventSet, 1),
		},
	}
}

// QueueEventRunner
//
// Takes a function literal(`func() func()`) that contains event execution content and sends it to the event processing channel.
func (ctx *ApplicationContext) QueueEventRunner(runner entity.EventRunner) {
	ctx.eventChannel.Channel <- runner
}

// QueueEventSet
//
// Sends an entity.EventSet with the entity publishing events
// and the entity.EventListener that receives and processes those events to the event distribution channel.
func (ctx *ApplicationContext) QueueEventSet(set entity.EventSet) {
	ctx.eventListenerDispenseChannel.DispenseChannel <- set
}

// ConsumeEventRunner
//
// Creates an event pool that processes events received from the channel.
// Event pools asynchronously process events.
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

		for i := 0; i < ctx.eventListenerDispenseChannel.DispensePoolSize; i++ {
			go func(innerContext context.Context) {
			selectLoop:
				for {
					select {
					case <-innerContext.Done():
						break selectLoop
					case eventSet := <-ctx.eventListenerDispenseChannel.DispenseChannel:
						manageEventRunnerContext(eventSet, func(set entity.EventSet) {
							ctx.eventChannel.Channel <- set.Runner
						})
					}
				}

				fmt.Println("[eventx] End of event listener dispense pool...")
			}(ctx.innerContext)
		}

		for i := 0; i < ctx.eventChannel.ProcessPoolSize; i++ {
			go func(innerContext context.Context) {
			selectLoop:
				for {
					select {
					case <-innerContext.Done():
						break selectLoop
					case runner := <-ctx.eventChannel.Channel:
						afterRunner := runner()
						if afterRunner != nil {
							ctx.eventChannel.AfterChannel <- afterRunner
						}
					}
				}

				fmt.Println("[eventx] End of event process pool...")
			}(ctx.innerContext)

			go func(innerContext context.Context) {
			selectLoop:
				for {
					select {
					case <-innerContext.Done():
						break selectLoop
					case runner := <-ctx.eventChannel.AfterChannel:
						runner()
					}
				}

				fmt.Println("[eventx] End of after event process pool...")
			}(ctx.innerContext)
		}
	})
}

// GetEventListener
//
// Returns the event listeners corresponding to the entity publishing the events.
// Since it can only return []any, type checking is required on the caller's side.
func (ctx *ApplicationContext) GetEventListener(typeVal reflect.Type) []any {
	listeners := ctx.eventListenerConfig.ListenerMap.Get(typeVal)
	return listeners
}

// RegisterEventListener
//
// Registers an event listener in the context.
func (ctx *ApplicationContext) RegisterEventListener(typeVal reflect.Type, eventListener any) {
	ctx.eventListenerConfig.ListenerMap.Put(typeVal, eventListener)
}

// Close
//
// Terminates the context, causing the event pool to end.
//
// Users are encouraged to execute this method with `defer eventx.Close` to ensure that `eventx` safely ends before their application terminates.
func (ctx *ApplicationContext) Close() {
	ctx.innerContextCancel()
}

// IsMultiMode
//
// Returns whether `eventx` can register and handle multiple event listeners for a single event entity.
func (ctx *ApplicationContext) IsMultiMode() bool {
	return ctx.eventListenerConfig.MultiEventMode
}
