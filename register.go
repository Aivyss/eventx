package eventx

import (
	"github.com/aivyss/eventx/errors"
	"reflect"
)

var listenerMap = map[reflect.Type]any{}

func RegisterEventListener[E any](el EventListener[E]) error {
	var e E
	typeVal := reflect.TypeOf(e)

	_, ok := listenerMap[typeVal]
	if ok {
		return errors.AlreadyRegisteredErr
	}

	listenerMap[typeVal] = el

	return nil
}

func Trigger[E any](elem E) error {
	typeVal := reflect.TypeOf(elem)
	listener, ok := listenerMap[typeVal]
	if !ok {
		return errors.NotFoundEventListenerErr
	}

	specifiedListener, ok := listener.(EventListener[E])
	if !ok {
		return errors.NotFoundEventListenerErr
	}

	go specifiedListener.Trigger(elem)

	return nil
}
