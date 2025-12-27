# Fluxo

Fluxo is a Go implementation of the **Event Bus pattern**, a messaging architecture that decouples event producers from event consumers.

## Features

- **Publish-Subscribe Model**: Decouple components through event-driven communication
- **Type-Safe Events**: Strongly-typed event handling
- **Async Event Processing**: Non-blocking event delivery
- **Multiple Subscribers**: Support for multiple handlers per event type
- **Simple API**: Minimal and intuitive interface

## Usage

### Basic Example

```go
import "github.com/mykytaserdiuk/fluxo/event"

// Create a bus
bus := event.New()

// Subscribe to events
bus.Subscribe("user:created", func(e Event) {
    fmt.Println("User created:", e.Data)
})

// Publish events
bus.Emit("user:created", Event{Data: "John"})
```
You can subscribe to events with different argument signatures:
```go
bus := event.New()

// Subscribe with args in handler
bus.Subscribe("new_follower", func(name string) {
    fmt.Println("NEW FOLLOWER: ", name)
})

// Subscribe without args
bus.Subscribe("new_follower", func() {
    fmt.Println("You have new follower!")
})

// Subscribe with multiple arguments
bus.Subscribe("new_follower", func(name, lastName string) { // or with pointer
    fmt.Println("Follower name: ", name, lastName)
})

// Emit event
bus.Emit("new_follower", "John")
```

See the [`/example`](/example/) directory for complete usage examples.

## Inspiration

This project is inspired by [EventBus](https://github.com/asaskevich/EventBus) and follows principles described in the [Event Bus Design Pattern](https://dzone.com/articles/design-patterns-event-bus).

<!-- ## License

MIT -->