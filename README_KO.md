# 목차
- [Installation](#installation)
- [개요](#개요)
- [애플리케이션 구동](#애플리케이션-구동)
- [EventListener](#eventlistener)
  * [인터페이스 설명](#인터페이스-설명)
  * [`EventListener[E]`의 임베디드 컴포지션](#eventlistenere의-임베디드-컴포지션)
  * [이벤트 등록](#이벤트-등록)
- [이벤트 triggering](#이벤트-triggering)
- [entity.EventContext](#entityeventcontext)
- [애플리케이션의 종료](#애플리케이션의-종료)

# Installation
```sh
go get github.com/aivyss/eventx@v1.3.0
```

# 개요
- eventx는 비동기 이벤트를 처리하기위한 애플리케이션입니다.
- 장점
  - 이벤트리스너를 이용해 각 로직계층간 의존성을 낮출수 있습니다.
  - 비동기로직 구성을 위한 별도의 복잡한 보일러플레이트 작성이 불필요해집니다.

# 애플리케이션 구동

- 애플리케이션은 자동으로 실행되지 않습니다. 당신의 애플리케이션 구동시 실행함수를 최소 1번은 실행하여야 합니다.
- 실행함수는 두가지가 있습니다.
  - `func RunDefaultApplication()`
  - `func RunApplication(eventChannelBufferSize int, eventProcessPoolSize int)`
  - 무엇을 고를지 애매하다면 `RunDefaultApplication`를 실행하길 권장합니다.

---

# EventListener

```go
type EventListener[E any] interface {
    Trigger(entity E) error
}
```

## 인터페이스 설명

- `E` 이벤트의 트리거가 될 대상입니다. 실행할 이벤트에 필요한 정보를 가져야 합니다.
  - `E` 구조체를 eventx애플리케이션에 전달시 이벤트가 트리거 됩니다.
  - `E` 는 사용자가 정의해야 합니다. (구조체, 인터페이스, 원시타입 등...)

- `Trigger(entity E) error` 이벤트를 처리하는 함수입니다.
- 당신은 이 인터페이스를 구현하고 등록해야 합니다.
- 인터페이스의 구현이 번거롭다면 아래의 메소드를 활용할 수 있습니다.
  ```go
  func RegisterFunAsEventListener[E any](trigger func(entity E) error) error
  
  func RegisterFuncThenAsEventListener[E any](
      trigger func(entity E) error,
      then func(entity E),
  ) error
  
  func RegisterFuncCatchAsEventListener[E any](
      trigger func(entity E) error,
      catch func(err error),
  ) error

  func RegisterFuncsAsEventListener[E any](
      trigger func(entity E) error,
      then func(entity E),
      catch func(err error),
  ) error
  ```
  - 당신은 이벤트가 트리거 될 때 작동할 lambda만 작성하면 됩니다.

<br>
<br>

## `EventListener[E]`의 임베디드 컴포지션
```go
type CatchErrEventListener[E any] interface {
    EventListener[E]
    Catch(err error)
}

type SuccessEventListener[E any] interface {
    EventListener[E]
    Then(entity E)
}

type CallbackEventListener[E any] interface {
    EventListener[E]
    SuccessEventListener[E]
    CatchErrEventListener[E]
}
```
- 이벤트가 실행된후 실행결과에 따라 후속으로 처리될 프로세스를 메소드로 포함하는 이벤트 리스너를 구현할 수 있습니다.

## 이벤트 등록

```go
func RegisterEventListener[E any](el EventListener[E]) error
```

- 당신이 이벤트 등록을 위해 사용할 기본적인 함수입니다.
- 이 함수로 등록되지 않은 이벤트리스너는 트리거되지 않습니다.
- 후속으로 보여드릴 편의성 Register함수들은 모두 내부적으로 이 함수를 호출하고 있습니다.

<br>
<br>

```go
func RegisterFuncAsEventListener[E any](trigger func(entity E) error) error
```

- `eventx` 는 이벤트등록의 편의성을 위해 `eventx.RegisterFuncAsEventListener`를 제공합니다.
- 이벤트리스너를 명시적인 별도의 코드나 변수로 유지시킬 필요가 없을 때, 사용하길 권장됩니다.

<br>
<br>



```go
func RegisterFuncCatchAsEventListener[E any](
  	trigger func(entity E) error,
  	catch func(err error),
) error

func RegisterFuncThenAsEventListener[E any](
  	trigger func(entity E) error,
  	then func(entity E),
) error

func RegisterFuncsAsEventListener[E any](
    trigger func(entity E) error,
    then func(entity E),
    catch func(err error),
) error
```

- 이 함수는 이벤트가 트리거 후의 후속 처리 절차도 함께 등록할 수 있습니다.
- 이벤트 트리거의 처리를 성공할 시, `then` 가 실행됩니다.
- 이벤트 트리거의 처리를 실패할 시, `catch`가 실행됩니다.
- 이벤트와 마찬가지로 후속 처리도 완전 비동기로 관리됩니다.
- 이벤트리스너를 명시적인 별도의 코드나 변수로 유지시킬 필요가 없을 때, 사용하길 권장됩니다.

# 이벤트 triggering

```go
func Trigger[E any](entity E) ([]entity.EventContext, error)
```

- 당신이 이벤트를 트리거 시키기 위해 사용해야할 함수는 단지 이것 뿐입니다.
- 이벤트를 트리거할 대상(entity)을 eventx 애플리케이션에 전달합니다.
- 전달된 대상은 비동기적으로 이벤트를 트리거하고 처리됩니다.
- `entity.EventContext`로 이벤트의 진행상황을 알 수 있으며 이벤트 실행전에 중단할 수도 있습니다.

# `entity.EventContext`
```go
type EventContext interface {
    IsRunnable() bool
    Cancel() bool
    IsDone() bool
}
```
- `IsRunnable`: `eventx`가 실행가능한 이벤트인지 여부를 반환합니다.
- `IsDone`: 이미 실행이 종료된 이벤트인지 여부를 반환합니다.
- `Cancel`: 실행전 이벤트라면 이벤트를 발행을 취소시킬 수 있습니다.

# 애플리케이션의 종료

```go
func Close()
```

- `Close` 함수를 호출해 `eventx` 애플리케이션을 종료할 수 있습니다.
- `defer` 키워드와 함께 활용해 당신의 애플리케이션을 종료하기 전에 `eventx` 를 먼저 종료시킬 수 있습니다.