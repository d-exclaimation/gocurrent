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
	. "github.com/d-exclaimation/gocurrent/types"
)

// Future is a data structure to represent suspended / deferred value
// while also adding additional utilities to work with that value.
//
// Example:
//
//  fut0 := future.Async(func() interface{} {
//      time.Sleep(10 * time.Second)
//      return 10, nil
//  })
//
type Future struct {
	// value is the internal wrapped value inside the Future
	value Any

	// err is the internal wrapped error value inside the Future
	err error

	// executor is an internal suspended function being used to hydrate value and err
	executor Function

	// Status is a value showing the status of the execution
	Status DeliveryStatus

	// _Delivery is a channel for the executor to send the finished result
	_Delivery chan result.Result
}

// Async instantiate a new Future and run the function in a separate goroutine.
func Async(exe Function) *Future {
	fut := &Future{
		value:     nil,
		err:       nil,
		executor:  exe,
		Status:    Idle,
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
	go func() {
		res := <-f._Delivery
		res.Match(result.Case{
			Success: func(value Any) {
				f.value = value
				f.Status = Success
			},
			Failure: func(err error) {
				f.err = err
				f.Status = Failure
			},
		})
		close(f._Delivery)
	}()
}

// Await waits for the Future to finish and return the values
//
// Note: Blocking operation!
func (f *Future) Await() (Any, error) {
	<-f.Done()
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

func (f *Future) Done() <-chan struct{} {
	await := make(chan struct{})
	go func() {
		for !f.IsDone() {
		}
		await <- struct{}{}
	}()
	return await
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
