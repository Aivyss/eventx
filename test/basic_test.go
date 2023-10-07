package test

import (
	"fmt"
	"github.com/aivyss/eventx"
	"testing"
)

func TestBasic(t *testing.T) {
	var l1 eventx.EventListener[string] = &Listener1{}
	var l2 eventx.EventListener[int] = &Listener2{}

	err := eventx.RegisterEventListener(l2)
	if err != nil {
		panic(err)
	}
	err = eventx.RegisterEventListener(l1)
	if err != nil {
		panic(err)
	}

	err = eventx.Trigger("test")
	if err != nil {
		panic(err)
	}
	err = eventx.Trigger(1)
	if err != nil {
		panic(err)
	}
}

type Listener1 struct{}

func (l *Listener1) Trigger(elem string) error {
	fmt.Println("Listener1 run =", elem)

	return nil
}

type Listener2 struct{}

func (l *Listener2) Trigger(elem int) error {
	fmt.Println("Listener2 run =", elem)
	return nil
}
