# Index
- [Installation](#installation)
- [Overview](#overview)
- [Running the Application](#running-the-application)
- [EventListener](#eventlistener)
  * [Interface Description](#interface-description)
  * [Embedded Composition of EventListener[E]](#embedded-composition-of-eventlistenere)
  * [Event Registration](#event-registration)
- [Event Triggering](#event-triggering)
- [entity.EventContext](#entityeventcontext)
- [Application Termination](#application-termination)

# Installation
```sh
go get github.com/aivyss/eventx@v1.3.0
```

# Overview
- eventx is an application for handling asynchronous events.
- Advantages
  - It allows lowering dependencies between logic layers using event listeners.
  - It eliminates the need for complex boilerplate code for configuring asynchronous logic.

<br>
<br>


# Running the Application
- The application does not automatically run. You must execute the run function at least once when starting your application.
- There are two run functions available:
  - `func RunDefaultApplication()`
  - `func RunApplication(eventChannelBufferSize int, eventProcessPoolSize int)`
  - If you are unsure, it is recommended to use `RunDefaultApplication`.

<br>
<br>

# EventListener
```go
type EventListener[E any] interface {
	Trigger(entity E) error
}
```

## Interface Description
- `E` represents the target that triggers the event. It should contain the necessary information for the event to be executed.
  - When you pass a E structure to the eventx application, the event is triggered.
  - `E` should be user-defined (a struct, interface, primitive type, etc.).
- `Trigger(entity E) error` is the function for processing events.
- You need to implement and register this interface.
- If implementing the interface is cumbersome, you can use the following method:
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
  - You only need to write a lambda that will operate when the event is triggered.
    <br>
    <br>

## Embedded Composition of EventListener[E]
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
- After an event is executed, you can implement an event listener that includes methods to process subsequent actions based on the execution result.

<br>
<br>

## Event Registration
```go
func RegisterEventListener[E any](el EventListener[E]) error
```
- This is the function you should use to register events.
- Event listeners that are not registered with this function will not be triggered.
- The convenience `Register` functions to be mentioned next all call this function internally.
<br>
<br>

```go
func RegisterFuncAsEventListener[E any](trigger func(entity E) error) error
```
- `eventx` provides `eventx.RegisterFuncAsEventListener` for the convenience of event registration.
- It is recommended for use when you do not need to maintain event listeners as explicit separate code or variables.
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
- This function allows you to register subsequent processing procedures after the event is triggered.
- When the event trigger processing is successful, `then` will be executed.
- When the event trigger processing fails, `catch` will be executed.
- Just like events, subsequent processing is also managed entirely asynchronously.
- It is recommended for use when you do not need to maintain event listeners as explicit separate code or variables.

<br>
<br>

# Event Triggering
```go
func Trigger[E any](elem E) ([]entity.EventContext, error)
```
- This is the function you should use to trigger events.
- Pass the target (elem) you want to trigger the event to the eventx application.
- The passed target is asynchronously triggered and processed.
- With `entity.EventContext`, you can track the progress of events and also stop them before execution.

# `entity.EventContext`
```go
type EventContext interface {
    IsRunnable() bool
    Cancel() bool
    IsDone() bool
}
```
- `IsRunnable`: Returns whether the event is executable by `eventx`.
- `IsDone`: Returns whether the event has already been executed.
- `Cancel`: If the event has not been executed yet, you can cancel the event publication.

# Application Termination
```go
func Close()
```
- You can use the `Close` function to terminate the eventx application.
- Utilize it with the `defer` keyword to ensure that `eventx` is terminated before your application exits.