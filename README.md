# Go-State: Lightweight & Reflect-Based State Machine

**Go-State** is a thread-safe, generic state machine for Go that leverages reflection to eliminate boilerplate. It allows you to dispatch events based on their types and manage state transitions with minimal friction.

## ✨ Features

* **Type-Safe Handlers:** Define handlers using specific event structs without manual type assertions.
* **Generic State:** Support for any custom state data structure using Go Generics (`[T any]`).
* **Parallel Composition:** Run multiple handlers for the same event concurrently with automatic context cancellation on failure.
* **Thread-Safe:** Built-in `sync.Mutex` to ensure consistent status transitions.
* **Context-Aware:** Full support for `context.Context` across all handlers.

## 📦 Installation

```bash
go get github.com/israelsodanoa/state
```

## 🚀 Quick Start

### 1. Define your State and Events
```go
type OrderState struct {
    ID int
}

type PaymentConfirmed struct {
    Amount float64
}
```

### 2. Setup the Machine
```go
func main() {
    ctx := context.Background()

    // Initialize with a state type and initial status
    sm := &state.StateMachine[OrderState]{
        State:  OrderState{ID: 42},
        Status: "PENDING",
    }

    // Create a handler for a specific event type
    confirmHandler := state.When(sm, func(ctx context.Context, e PaymentConfirmed) error {
        fmt.Printf("Processing payment of $%.2f\n", e.Amount)
        return nil
    })

    // Map the handler to a status transition
    sm.TransitionTo("COMPLETED", confirmHandler)

    // 3. Publish an event
    // The library automatically finds the handler for the 'PaymentConfirmed' type
    err := sm.Pub(ctx, PaymentConfirmed{Amount: 150.00})
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("New Status:", sm.Status) // Output: New Status: COMPLETED
}
```

---

## 🛠 Advanced Usage

### Parallel Handler Composition
One of the most powerful features of **Go-SM** is the `Compose` function. It allows you to trigger multiple side effects (like logging, analytics, or notifications) for a single event in parallel.

```go
logHandler := state.When(sm, func(ctx context.Context, e PaymentConfirmed) error {
    log.Println("Logging to database...")
    return nil
})

emailHandler := state.When(sm, func(ctx context.Context, e PaymentConfirmed) error {
    log.Println("Sending confirmation email...")
    return nil
})

// Compose them! They will run in separate goroutines.
// If one fails, the context for the others is cancelled.
combined := state.Compose(logHandler, emailHandler)

sm.TransitionTo("READY_FOR_SHIPPING", combined)
```

---

## 📚 API Overview

### `StateMachine[T]`
The core engine. It holds your custom state `T` and a slice of `EventTransition`.
* `Pub(ctx, data)`: Matches the type of `data` to a registered handler and executes it.
* `TransitionTo(status, handler)`: Registers which status the machine should move to after a handler successfully executes.

### `When[T, E](sm, handlerFn)`
A helper to register a handler function.
* **T**: Your state type.
* **E**: The event struct type you want to listen for.
* *Safety:* This will `panic` if you try to register two different handlers for the same event type on the same machine.

### `Compose(handlers...)`
Takes multiple `EventHandler` instances and returns a single wrapped handler. 
* **Execution:** All handlers run in their own goroutine.
* **Error Handling:** Returns the first error encountered and cancels the execution of the remaining handlers in that group.

---

## ⚖️ License
Distributed under the MIT License. See `LICENSE` for more information.
