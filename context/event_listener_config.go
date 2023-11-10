package context

import (
	"github.com/aivyss/typex"
	"reflect"
)

type EventListenerConfig struct {
	MultiEventMode bool
	ListenerMap    typex.MultiMap[reflect.Type, any]
}
