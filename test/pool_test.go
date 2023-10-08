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
}

func TestPPool2(t *testing.T) {
	t.Run("regenerate application context", func(t *testing.T) {
		eventx.RunDefaultApplication()
		time.Sleep(500 * time.Millisecond)
		defaultGoroutineNum := runtime.NumGoroutine()

		eventx.RunApplication(3, 15)

		for {
			time.Sleep(500 * time.Millisecond)

			if defaultGoroutineNum+5 == runtime.NumGoroutine() {
				fmt.Println("[success] regenerate application context")
				eventx.Close()
				break
			}
		}
	})
}

func TestPool3(t *testing.T) {
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

func TestPool4(t *testing.T) {
	t.Run("async test", func(t *testing.T) {
		var mutex sync.Mutex
		flag := false
		listener3Count := 0
		loopCnt := 10
		eventx.RunDefaultApplication()
		defer eventx.Close()

		_ = eventx.RegisterFuncAsEventListener(func(entity Listener3Elem) error {
			mutex.Lock()
			listener3Count += 1
			mutex.Unlock()

			time.Sleep(1 * time.Second)

			mutex.Lock()
			flag = true
			mutex.Unlock()

			return nil
		})
		for i := 0; i < loopCnt; i++ {
			_ = eventx.Trigger(Listener3Elem(i))
		}

		time.Sleep(100 * time.Millisecond)

		if flag || listener3Count != loopCnt {
			panic("[fail] async test")
		}

		for {
			time.Sleep(500 * time.Millisecond)

			if flag {
				fmt.Println("[success] async test")
				break
			}
		}
	})
}
