package test

import (
	"errors"
	"fmt"
	"github.com/aivyss/eventx"
	"github.com/aivyss/eventx/entity"
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
			_, _ = eventx.Trigger(TestEventEntity(i))
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
		defer eventx.Close()

		for {
			time.Sleep(500 * time.Millisecond)

			if defaultGoroutineNum+10 == runtime.NumGoroutine() {
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

		listener := entity.BuildEventListener(func(entity TestEventEntity) error {
			time.Sleep(1 * time.Millisecond)
			mutex.Lock()
			listener3Count += 1
			mutex.Unlock()

			return nil
		})

		_ = eventx.RegisterEventListener(listener)

		for i := 0; i < loopCnt; i++ {
			_, _ = eventx.Trigger(TestEventEntity(i))
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
			_, _ = eventx.Trigger(TestEventEntity(i))
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
		defer eventx.Close()

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
			_, _ = eventx.Trigger(TestEventEntity(i))
			_, _ = eventx.Trigger(TestEventEntity2(i))
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

func TestCallbackPool(t *testing.T) {
	t.Run("callback test", func(t *testing.T) {
		var mutex sync.Mutex
		listener3Count := 0
		loopCnt := 10
		eventx.RunDefaultApplication()
		defer eventx.Close()

		eventx.RegisterFuncsAsEventListener(
			func(entity TestEventEntity) error {
				return errors.New("test_error =" + fmt.Sprint(entity))
			},
			func(entity TestEventEntity) {
				panic("DO NOT RUN")
			},
			func(err error) {
				if err == nil {
					panic("err is nil")
				}

				mutex.Lock()
				listener3Count += 1
				mutex.Unlock()
			},
		)
		eventx.RegisterFuncsAsEventListener(
			func(entity TestEventEntity) error {
				return nil
			},
			func(entity TestEventEntity) {
				mutex.Lock()
				listener3Count += 1
				mutex.Unlock()
			},
			func(err error) {
				panic("DO NOT RUN")
			},
		)
		eventx.RegisterFuncsAsEventListener(
			func(entity TestEventEntity) error {
				return errors.New("test_error =" + fmt.Sprint(entity))
			},
			nil,
			func(err error) {
				if err == nil {
					panic("err is nil")
				}

				mutex.Lock()
				listener3Count += 1
				mutex.Unlock()
			},
		)
		eventx.RegisterFuncsAsEventListener(
			func(entity TestEventEntity) error {
				return nil
			},
			func(entity TestEventEntity) {
				mutex.Lock()
				listener3Count += 1
				mutex.Unlock()
			},
			nil,
		)

		for i := 0; i < loopCnt; i++ {
			eventx.Trigger(TestEventEntity(i))
		}

		for {
			time.Sleep(500 * time.Millisecond)

			if listener3Count == loopCnt*4 {
				fmt.Println("[success] callback test")
				break
			}
		}
	})
}

func TestTemp(t *testing.T) {
	flag := false
	eventx.RunDefaultApplication()
	defer eventx.Close()

	_ = eventx.RegisterFuncAsEventListener(func(entity TestEventEntity) error {
		flag = true
		return nil
	})

	eventCtxs, _ := eventx.Trigger(TestEventEntity(1))
	eventCtx := eventCtxs[0]

	for {
		time.Sleep(500 * time.Millisecond)

		if flag && !eventCtx.IsRunnable() && eventCtx.IsDone() {
			fmt.Println("[success] cancel functions - 1")
			break
		}
	}
}

func TestTemp2(t *testing.T) {
	flag := false
	eventx.RunDefaultApplication()
	defer eventx.Close()

	_ = eventx.RegisterFuncAsEventListener(func(entity TestEventEntity) error {
		flag = true
		return nil
	})

	eventCtxs, _ := eventx.Trigger(TestEventEntity(1))
	eventCtx := eventCtxs[0]
	eventCtx.Cancel()

	for {
		time.Sleep(500 * time.Millisecond)

		if !flag && !eventCtx.IsRunnable() && !eventCtx.IsDone() {
			fmt.Println("[success] cancel function - 2")
			break
		}
	}
}
