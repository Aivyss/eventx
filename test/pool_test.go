package test

import (
	"fmt"
	"github.com/aivyss/eventx"
	"runtime"
	"sync"
	"testing"
	"time"
)

type TestEventEntity int

func TestPool(t *testing.T) {
	t.Run("RegisterFuncAsEventListener", func(t *testing.T) {
		var mutex sync.Mutex
		listener3Count := 0
		loopCnt := 1000

		eventx.RunDefaultApplication()
		defer eventx.Close()

		_ = eventx.RegisterFuncAsEventListener(func(entity TestEventEntity) error {
			mutex.Lock()
			listener3Count += 1
			mutex.Unlock()

			return nil
		})

		for i := 0; i < loopCnt; i++ {
			_ = eventx.Trigger(TestEventEntity(i))
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

		eventx.RunApplication(3, 15, true)

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

		listener := eventx.BuildEventListener(func(entity TestEventEntity) error {
			time.Sleep(1 * time.Millisecond)
			mutex.Lock()
			listener3Count += 1
			mutex.Unlock()

			return nil
		})

		_ = eventx.RegisterEventListener(listener)

		for i := 0; i < loopCnt; i++ {
			_ = eventx.Trigger(TestEventEntity(i))
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

		_ = eventx.RegisterFuncAsEventListener(func(entity TestEventEntity) error {
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
			_ = eventx.Trigger(TestEventEntity(i))
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

func TestPool5(t *testing.T) {
	var mutex sync.Mutex
	event1Cnt := 0
	event2Cnt := 0
	event3Cnt := 0
	eventTriggeredCnt := 0
	loopCnt := 100

	t.Run("multi-test", func(t *testing.T) {
		eventx.RunDefaultApplication()

		type TestEventEntity2 int
		_ = eventx.RegisterFuncAsEventListener(func(entity TestEventEntity2) error {
			mutex.Lock()
			event3Cnt += 1
			mutex.Unlock()
			time.Sleep(5 * time.Millisecond)

			return nil
		})
		_ = eventx.RegisterFuncAsEventListener(func(entity TestEventEntity) error {
			mutex.Lock()
			eventTriggeredCnt += 1
			event1Cnt += 1
			mutex.Unlock()
			time.Sleep(20 * time.Millisecond)

			return nil
		})
		_ = eventx.RegisterFuncAsEventListener(func(entity TestEventEntity) error {
			mutex.Lock()
			eventTriggeredCnt += 1
			event2Cnt += 1
			mutex.Unlock()
			time.Sleep(20 * time.Millisecond)

			return nil
		})

		for i := 0; i < loopCnt; i++ {
			_ = eventx.Trigger(TestEventEntity(i))
			_ = eventx.Trigger(TestEventEntity2(i))
		}

		for {
			time.Sleep(500 * time.Millisecond)

			if event1Cnt == loopCnt && event2Cnt == loopCnt && eventTriggeredCnt == loopCnt*2 && event3Cnt == loopCnt {
				fmt.Println("[success] multi-test")
				break
			}
		}
	})
}
