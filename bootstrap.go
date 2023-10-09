package eventx

import (
	"github.com/aivyss/eventx/entity"
	"github.com/aivyss/eventx/errors"
	"reflect"
)

func RunDefaultApplication() {
	RunApplication(entity.DefaultEventChannelBufferSize, entity.DefaultEventProcessPoolSize, entity.DefaultMultiEventMode)
}

func RunApplication(eventChannelBufferSize int, eventProcessPoolSize int, multiEventMode bool) {
	if appContext != nil {
		appContext.Close()
	}

	appContext = entity.NewApplicationContext(eventChannelBufferSize, eventProcessPoolSize, multiEventMode)
	appContext.ConsumeEventRunner()
}

func RegisterEventListener[E any](el EventListener[E]) error {
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
	return RegisterEventListener(BuildEventListener(trigger))
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

	var specifiedListeners []EventListener[E]
	for _, listener := range listeners {
		specifiedListener, ok := listener.(EventListener[E])
		if !ok {
			return errors.NotFoundEventListenerErr
		}

		specifiedListeners = append(specifiedListeners, specifiedListener)
	}

	generateEventRunner := func(elem E, listener EventListener[E]) entity.EventRunner {
		return func() {
			_ = listener.Trigger(elem)
		}
	}

	for _, listener := range specifiedListeners {
		appContext.QueueEventRunner(generateEventRunner(elem, listener))
	}

	return nil
}
