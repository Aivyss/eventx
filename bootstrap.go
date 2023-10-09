package eventx

import (
	"github.com/aivyss/eventx/context"
	"github.com/aivyss/eventx/entity"
	"github.com/aivyss/eventx/errors"
	"reflect"
)

func RunDefaultApplication() {
	RunApplication(context.DefaultEventChannelBufferSize, context.DefaultEventProcessPoolSize, context.DefaultMultiEventMode)
}

func RunApplication(eventChannelBufferSize int, eventProcessPoolSize int, multiEventMode bool) {
	if appContext != nil {
		appContext.Close()
	}

	appContext = context.NewApplicationContext(eventChannelBufferSize, eventProcessPoolSize, multiEventMode)
	appContext.ConsumeEventRunner()
}

func RegisterEventListener[E any](el entity.EventListener[E]) error {
	var e E
	typeVal := reflect.TypeOf(e)

	listeners := appContext.GetEventListener(typeVal)
	if !appContext.IsMultiMode() && len(listeners) > 0 {
		return errors.AlreadyRegisteredErr
	}

	appContext.RegisterEventListener(typeVal, el)

	return nil
}

func RegisterFuncAsEventListener[E any](trigger func(entity E) error) error {
	return RegisterEventListener(entity.BuildEventListener(trigger))
}

func Close() {
	appContext.Close()
}

func Trigger[E any](elem E) error {
	typeVal := reflect.TypeOf(elem)
	listeners := appContext.GetEventListener(typeVal)
	if len(listeners) == 0 {
		return errors.NotFoundEventListenerErr
	}

	for _, listener := range listeners {
		specifiedListener, ok := listener.(entity.EventListener[E])
		if !ok {
			return errors.NotFoundEventListenerErr
		}

		appContext.QueueEventSet(entity.NewEventSet(specifiedListener, elem))
	}

	return nil
}
