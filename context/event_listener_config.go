package context

import (
	"github.com/aivyss/eventx/common"
	"reflect"
)

type EventListenerConfig struct {
	MultiEventMode bool
	ListenerMap    common.MultiMap[reflect.Type, any]
}
