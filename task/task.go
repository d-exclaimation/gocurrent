//
//  task.go
//  task
//
//  Created by d-exclaimation on 6:38 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package task

import (
	"github.com/d-exclaimation/gocurrent/try"
	"time"
)

// Task is  unit of asynchronous work
//
// Example code:
//    t0 := task.Async[int](func() (int, error) {
//        time.Sleep(1 * time.Second)
//        return 10, nil
//    }
//    task.AsyncVoid(func() error {
//        res, err := t0.Await()
//        if err != nil {
//            fmt.Println(err.Error())
//        } else {
//            fmt.Println("We got " + res)
//        }
//        return nil
//    }
type Task[T any] struct {
	// Wrapped value in a Try of the Task
	value *try.Try[T]
	// The function to acquire the value or an error
	process func() (T, error)
	// The channels that requested for the awaited value
	deliveries map[chan<- *try.Try[T]]<-chan *try.Try[T]
	// Channel for sending newly acquired value
	mailbox chan *try.Try[T]
	// Channel for sending new awaiter
	delivery chan chan *try.Try[T]
}

// New creates a new Task but does not run it.
//
// Running will be delegated to the caller by running task.Run or task.LazyAwait.
func New[T any](op func() (T, error)) *Task[T] {
	task := &Task[T]{
		value:      nil,
		process:    op,
		mailbox:    make(chan *try.Try[T]),
		deliveries: make(map[chan<- *try.Try[T]]<-chan *try.Try[T]),
		delivery:   make(chan chan *try.Try[T]),
	}
	task.behavior()
	return task
}

// Async creates a new Task and run it immediately.
func Async[T any](op func() (T, error)) *Task[T] {
	task := New[T](op)
	task.Run()
	return task
}

// AsyncVoid run a non-returning function in a Task
func AsyncVoid(op func() error) {
	New[struct{}](func() (struct{}, error) {
		return struct{}{}, op()
	}).Run()
}

// Run the Task, if already run before it will retry running but does not reset state
// i.e. the previous value will not be cleared until the new value acquired.
func (t *Task[T]) Run() {
	go func() {
		t.mailbox <- try.From[T](t.process)
	}()
}

// Actor like behavior
func (t *Task[T]) behavior() {
	go func() {
		for {
			select {
			case deliver, valid := <-t.delivery:
				if !valid {
					continue
				}
				if t.value != nil {
					deliver <- t.value
					close(deliver)
				} else {
					t.deliveries[deliver] = deliver
				}

			case res, valid := <-t.mailbox:
				if !valid {
					continue
				}
				t.value = res
				if res != nil {
					for in, _ := range t.deliveries {
						in <- res
						close(in)
						delete(t.deliveries, in)
					}
				}
			}
		}
	}()
}

// LazyAwait run the Task and waits for the returned value.
// Important:
//  - This will not work for retries
//  - Use only after task.New
func (t *Task[T]) LazyAwait() (T, error) {
	t.Run()
	return t.Await()
}

// Await for the asynchronous value.
func (t *Task[T]) Await() (T, error) {
	res := t.Try()
	return res.ToOption()
}

// AwaitChannel returns the consuming channel for awaiting the asynchronous value.
func (t *Task[T]) AwaitChannel() <-chan *try.Try[T] {
	deliver := make(chan *try.Try[T])
	go func() {
		time.Sleep(time.Millisecond)
		t.delivery <- deliver
	}()
	return deliver
}

// Try awaits for the asynchronous value but return them in a Try.
func (t *Task[T]) Try() *try.Try[T] {
	deliver := t.AwaitChannel()
	return <-deliver
}

// Then perform callback on value acquired.
func (t *Task[T]) Then(matcher try.Case[T]) {
	go func() {
		t.Try().Match(matcher)
	}()
}
