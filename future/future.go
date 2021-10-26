//
//  future.go
//  future
//
//  Created by d-exclaimation on 1:17 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package future

import (
	"errors"
	"github.com/d-exclaimation/gocurrent/result"
)

// Future is a data structure to represent suspended / deferred value
// while also adding additional utilities to work with that value
type Future struct {
	// value is the internal wrapped value inside the Future
	value Any

	// err is the internal wrapped error value inside the Future
	err error

	// executor is an internal suspended function being used to hydrate value and err
	executor Function

	// Status is a value showing the status of the execution
	Status DeliveryStatus

	// awaiter is a store of other channels waiting for the finish of the execution
	awaiter map[chan bool]bool

	// _Delivery is a channel for the executor to send the finished result
	_Delivery chan result.Result

	// _Awaiter is a channel to register a new waiting channel in awaiter
	_Awaiter chan chan bool
}

// Async instantiate a new Future and run the function in a separate goroutine.
func Async(exe Function) *Future {
	fut := &Future{
		value:     nil,
		err:       nil,
		executor:  exe,
		Status:    Idle,
		awaiter:   make(map[chan bool]bool),
		_Awaiter:  make(chan chan bool),
		_Delivery: make(chan result.Result),
	}
	fut.Reload()
	return fut
}

// New instantiate a new Future but does not run the function.
func New(exe Function) *Future {
	return &Future{
		value:     nil,
		err:       nil,
		executor:  exe,
		Status:    Idle,
		awaiter:   make(map[chan bool]bool),
		_Awaiter:  make(chan chan bool),
		_Delivery: make(chan result.Result),
	}
}

// Reload run the actor-like receiver and the future's function in separate goroutines
// to hydrate the value and error
func (f *Future) Reload() {
	f.behavior()
	f.run()
}

// run execute the function is a separate goroutine, set the status, and pipe back the result
func (f *Future) run() {
	f.Status = Loading
	go func() {
		data, err := f.executor()
		f._Delivery <- result.New(data, err)
	}()
}

// behavior execute the receiver on a separate goroutine.
func (f *Future) behavior() {
	go f.receive()
}

// receive create an actor-like behavior
// to concurrent-safely queue-in  messages from the channels
func (f *Future) receive() {
	for {
		select {

		case res := <-f._Delivery:
			res.Match(result.Case{
				Success: func(value interface{}) {
					f.value = value
					f.Status = Success
				},
				Failure: func(err error) {
					f.err = err
					f.Status = Failure
				},
			})

			// Notify waiters
			for ch := range f.awaiter {
				ch <- true
			}
			return
		case ch := <-f._Awaiter:
			f.awaiter[ch] = true
		}
	}
}

// Await waits for the Future to finish and return the values
//
// Note: Blocking operation!
func (f *Future) Await() (Any, error) {
	if !f.IsDone() {
		await := make(chan bool)
		f._Awaiter <- await
		<-await
		close(await)
	}
	return f.Get()
}

// Get directly try to access the internal values if possible
func (f *Future) Get() (Any, error) {
	if !f.IsDone() {
		return nil, errors.New("future 'Get': future has not concluded, values are missing, try using 'Await()' instead")
	}
	return f.value, f.err
}

// Result act similar to Get but return a result.Result instead
func (f *Future) Result() result.Result {
	data, err := f.Get()
	return result.New(data, err)
}

// IsSuccess indicate whether the Future has completed and return a successful value
func (f *Future) IsSuccess() bool {
	return f.Status == Success
}

// IsDone indicate whether the Future has completed
func (f *Future) IsDone() bool {
	return f.Status == Success || f.Status == Failure
}

// Channel convert the result into a single consumer-only channel
func (f *Future) Channel() <-chan result.Result {
	ch := make(chan result.Result)
	go func() {
		data, err := f.Await()
		res := result.New(data, err)
		ch <- res
		close(ch)
	}()
	return ch
}
