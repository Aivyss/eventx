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
	if trigger == nil {
		return errors.NoTriggerFuncErr
	}

	return RegisterEventListener(entity.BuildEventListener(trigger))
}

func RegisterFuncThenAsEventListener[E any](
	trigger func(entity E) error,
	then func(entity E),
) error {
	if trigger == nil {
		return errors.NoTriggerFuncErr
	}
	if then == nil {
		return RegisterFuncAsEventListener(trigger)
	}

	return RegisterEventListener(entity.BuildSuccessEventListener(trigger, then))
}

func RegisterFuncCatchAsEventListener[E any](
	trigger func(entity E) error,
	catch func(err error),
) error {
	if trigger == nil {
		return errors.NoTriggerFuncErr
	}
	if catch == nil {
		return RegisterFuncAsEventListener(trigger)
	}

	return RegisterEventListener(entity.BuildCatchErrEventListener(trigger, catch))
}

func RegisterFuncsAsEventListener[E any](
	trigger func(entity E) error,
	then func(entity E),
	catch func(err error),
) error {
	if trigger == nil {
		return errors.NoTriggerFuncErr
	}
	if then == nil && catch == nil {
		return RegisterFuncAsEventListener(trigger)
	}
	if then == nil {
		return RegisterFuncCatchAsEventListener(trigger, catch)
	}
	if catch == nil {
		return RegisterFuncThenAsEventListener(trigger, then)
	}

	return RegisterEventListener(entity.BuildEventListenerWithCallback(
		trigger,
		then,
		catch,
	))
}

func Close() {
	appContext.Close()
}

func Trigger[E any](elem E) ([]entity.EventContext, error) {
	typeVal := reflect.TypeOf(elem)
	listeners := appContext.GetEventListener(typeVal)
	if len(listeners) == 0 {
		return nil, errors.NotFoundEventListenerErr
	}

	var ctxs []entity.EventContext
	for _, listener := range listeners {
		specifiedListener, ok := listener.(entity.EventListener[E])
		if !ok {
			return nil, errors.NotFoundEventListenerErr
		}

		set := entity.NewEventSet(specifiedListener, elem)
		ctxs = append(ctxs, set.Context())
		appContext.QueueEventSet(set)
	}

	return ctxs, nil
}
