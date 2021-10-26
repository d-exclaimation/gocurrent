//
//  jet.go
//  jet
//
//  Created by d-exclaimation on 7:21 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package jet

import "errors"

// Jet is a data structure for streaming like behavior with a singular upstream and multiple consumer channel.
//
//  jt := jet.New()
//  for value := range jt.Sink() {
//      log.Println(value)
//  }
//
// with multiple utilities for handling time-based value.
//
//  jt := jet.New()
//  jt.OnSnapshot(func(value interface{}) {
//      log.Println(value)
//  })
//
// Also an iterator that can iterate using the Next, Value, and Err method in a for loop (blocking).
//
//  jt := jet.New()
//  for jt.Next() {
//      log.Println(jt.Value())
//  }
//
// Also handle single recent value request with caching and provide method like Await and AwaitNoCache.
// Also handle closing all channels and deallocating resources.
type Jet struct {
	// _upstream is the upstream channel to push data into the Jet
	_upstream chan interface{}

	// _register is the channel to concurrently set a new consumer channel
	_register chan chan interface{}

	// _unregister is the channel to concurrently unset and close a consumer channel
	_unregister chan Consumer

	// _await is the channel for sending single use channel
	_await chan chan interface{}

	// _stop is the shutdown channel
	_stop chan bool

	// cache is the preserved latest value
	cache interface{}

	// err is the accumulated errors
	err error

	// downstream is the map state for store long-running consumer to producer channel pair
	downstream Downstreams

	// waiters is the map state for store single use channel
	waiters OneTimers

	// isDone is the state to indicate whether Jet finished
	isDone bool
}

// New instantiate a new Jet and run the behavior in a separate goroutine.
func New() *Jet {
	j := &Jet{
		_upstream:   make(chan interface{}),
		_register:   make(chan chan interface{}),
		_unregister: make(chan Consumer),
		_await:      make(chan chan interface{}),
		_stop:       make(chan bool),
		cache:       nil,
		err:         nil,
		downstream:  make(Downstreams),
		waiters:     make(OneTimers),
		isDone:      false,
	}
	j.behavior()
	return j
}

// behavior is a method for running the receiver
func (j *Jet) behavior() {
	go j.receive()
}

// receive is method for actor-like behavior for handling messages from channels
func (j *Jet) receive() {
	for {
		select {
		// Up the value to all consumer and close all awaiter
		case elem := <-j._upstream:
			j.emit(elem)

		// Register a consumer and unregister one
		case channel := <-j._register:
			j.downstream[channel] = channel
		case consumer := <-j._unregister:
			producer, ok := j.downstream[consumer]
			if ok {
				close(producer)
				delete(j.downstream, consumer)
			}

		// Single value consumer
		case await := <-j._await:
			j.waiters[await] = true

		case stop := <-j._stop:
			j.isDone = stop
			if stop {
				j.shutdown()
				return
			}
		}
	}
}

// emit dispatch all the element to all downstream and waiters
func (j *Jet) emit(elem interface{}) {
	j.cache = elem
	for _, producer := range j.downstream {
		producer <- elem
	}
	for awaiter := range j.waiters {
		awaiter <- elem
		close(awaiter)
		delete(j.waiters, awaiter)
	}
}

// shutdown close all downstream, waiters, and channels
func (j *Jet) shutdown() {
	for consumer, producer := range j.downstream {
		close(producer)
		delete(j.downstream, consumer)
	}
	for awaiter := range j.waiters {
		awaiter <- j.cache
		close(awaiter)
		delete(j.waiters, awaiter)
	}
	close(j._await)
	close(j._register)
	close(j._unregister)
	close(j._stop)
}

// Up pushes a new value into the Jet
func (j *Jet) Up(data interface{}) {
	if j.isDone {
		return
	}

	j._upstream <- data
}

// Close shutdown the entire Jet and all downstream from Sink
func (j *Jet) Close() {
	if j.isDone {
		return
	}

	j._stop <- true
	defer close(j._upstream)
}

// Sink registers a consumer channel and return it
func (j *Jet) Sink() Consumer {
	consumer := make(chan interface{})

	if j.isDone {
		defer close(consumer)
	} else {
		j._register <- consumer
	}

	return consumer
}

// Unsink unregisters a consumer channel and return an error
func (j *Jet) Unsink(ch Consumer) error {
	if j.isDone {
		return errors.New("jet 'Unlink': Jet has finished or been shutdown forcefully")
	}
	j._unregister <- ch
	return nil
}

// OnSnapshot register sink and iterate over it and call the callback
func (j *Jet) OnSnapshot(callback func(interface{})) (<-chan bool, func()) {
	ch := make(chan interface{})
	done := make(chan bool)

	if j.isDone {
		defer close(ch)
	} else {
		j._register <- ch
	}

	go func() {
		for elem := range ch {
			callback(elem)
		}
		done <- true
	}()

	return done, func() {
		if j.isDone {
			return
		}
		j._unregister <- ch
	}
}

// Await is method for waiting for the next value in the Jet otherwise use the cache
func (j *Jet) Await() interface{} {
	res := j.AwaitNoCache()
	if res != nil {
		return j.cache
	}
	return res
}

// AwaitNoCache is a method for waiting for the next value in the Jet but doesn't use the cache
func (j *Jet) AwaitNoCache() interface{} {
	if j.isDone {
		return nil
	}

	await := make(chan interface{})
	j._await <- await
	return <-await
}

// --- Iterator ---

// Next give back a boolean to indicate whether the iterator finished
func (j *Jet) Next() bool {
	res := j.AwaitNoCache()
	return res != nil || !j.isDone
}

// Value return the current value in the iteration
//
// Note: To get next value, call Next method
func (j *Jet) Value() interface{} {
	return j.cache
}

// Err return the accumulated error from the Jet iterator
func (j *Jet) Err() error {
	return j.err
}
