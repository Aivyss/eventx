package test

import (
	"fmt"
	"github.com/aivyss/eventx"
	"runtime"
	"sync"
	"testing"
	"time"
)

type Listener3Elem int

func TestPool(t *testing.T) {
	t.Run("RegisterFuncAsEventListener", func(t *testing.T) {
		var mutex sync.Mutex
		listener3Count := 0
		loopCnt := 1000

		eventx.RunDefaultApplication()
		defer eventx.Close()

		_ = eventx.RegisterFuncAsEventListener(func(entity Listener3Elem) error {
			mutex.Lock()
			listener3Count += 1
			mutex.Unlock()

			return nil
		})

		for i := 0; i < loopCnt; i++ {
			_ = eventx.Trigger(Listener3Elem(i))
		}

		for {
			time.Sleep(500 * time.Millisecond)

			if listener3Count == loopCnt {
				fmt.Println("[success] RegisterFuncAsEventListener")
				break
			}
		}
	})

	t.Run("regenerate application context", func(t *testing.T) {
		eventx.RunDefaultApplication()
		defaultGoroutineNum := runtime.NumGoroutine()

		eventx.RunApplication(3, 15)
		defer eventx.Close()

		for {
			time.Sleep(500 * time.Millisecond)

			if defaultGoroutineNum+5 == runtime.NumGoroutine() {
				fmt.Println("[success] regenerate application context")
				break
			}
		}
	})

	t.Run("event test", func(t *testing.T) {
		var mutex sync.Mutex
		listener3Count := 0
		loopCnt := 1000

		eventx.RunDefaultApplication()
		defer eventx.Close()

		listener := eventx.BuildEventListener(func(entity Listener3Elem) error {
			time.Sleep(1 * time.Millisecond)
			mutex.Lock()
			listener3Count += 1
			mutex.Unlock()

			return nil
		})

		_ = eventx.RegisterEventListener(listener)

		for i := 0; i < loopCnt; i++ {
			_ = eventx.Trigger(Listener3Elem(i))
		}

		for {
			if listener3Count == loopCnt {
				fmt.Println("[success] event test")
				break
			}
			time.Sleep(5 * time.Millisecond)
		}
	})
}
