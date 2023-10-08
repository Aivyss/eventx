package eventx

import (
	"github.com/aivyss/eventx/entity"
	"github.com/aivyss/eventx/errors"
	"reflect"
)

func RunDefaultApplication() {
	RunApplication(entity.DefaultEventChannelBufferSize, entity.DefaultEventProcessPoolSize)
}

func RunApplication(eventChannelBufferSize int, eventProcessPoolSize int) {
	if appContext != nil {
		appContext.Close()
	}

	appContext = entity.NewApplicationContext(eventChannelBufferSize, eventProcessPoolSize)
	appContext.ConsumeEventRunner()
}

func RegisterEventListener[E any](el EventListener[E]) error {
	var e E
	typeVal := reflect.TypeOf(e)

	_, ok := appContext.GetEventListener(typeVal)
	if ok {
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
	listener, ok := appContext.GetEventListener(typeVal)
	if !ok {
		return errors.NotFoundEventListenerErr
	}

	specifiedListener, ok := listener.(EventListener[E])
	if !ok {
		return errors.NotFoundEventListenerErr
	}

	runner := func() {
		_ = specifiedListener.Trigger(elem)
	}
	appContext.QueueEventRunner(runner)

	return nil
}
