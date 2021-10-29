//
//  jet.go
//  jet
//
//  Created by d-exclaimation on 7:21 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package jet

import (
	"context"
	"errors"
	"github.com/d-exclaimation/gocurrent/streaming"
	. "github.com/d-exclaimation/gocurrent/types"
	"log"
)

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
//  jt.On(func(value Any) {
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
	// upstream is the upstream channel to push data into the Jet
	upstream chan Any

	// registrar is the channel to concurrently set a new consumer channel
	registrar chan chan Any

	// unregistrar is the channel to concurrently unset and close a consumer channel
	unregistrar chan streaming.Consumer

	// awaiter is the channel for sending single use channel
	awaiter chan chan Any

	// acid is the shutdown channel
	acid chan Signal

	// latestSnapshot is the preserved latest value
	latestSnapshot Any

	// accumulatedError is the accumulated errors
	accumulatedError error

	// downstream is the map state for store long-running consumer to producer channel pair
	downstream streaming.Downstreams

	// waiters is the map state for store single use channel
	waiters streaming.Downstreams

	// isDone is the state to indicate whether Jet finished
	isDone bool
}

// New instantiate a new Jet and run the behavior in a separate goroutine.
func New(opts ...Option) *Jet {
	return Lazy(opts...)()
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
		case snapshot, valid := <-j.upstream:
			if !valid {
				continue
			}
			j.emit(snapshot)

		// Register a consumer and unregister one
		case channel, valid := <-j.registrar:
			if !valid {
				continue
			}
			j.downstream[channel] = channel
		case consumer, valid := <-j.unregistrar:
			if !valid {
				continue
			}
			producer, ok := j.downstream[consumer]
			if ok {
				close(producer)
				delete(j.downstream, consumer)
			}

		// Single value consumer
		case await, valid := <-j.awaiter:
			if !valid {
				continue
			}
			j.waiters[await] = await

		case _, valid := <-j.acid:
			if !valid {
				continue
			}
			j.isDone = true
			j.shutdown()
			return
		}
	}
}

// emit dispatch all the element to all downstream and waiters
func (j *Jet) emit(snapshot Any) {
	j.latestSnapshot = snapshot
	for _, producer := range j.downstream {
		producer <- snapshot
	}
	for awaitConsumer, awaitProducer := range j.waiters {
		awaitProducer <- snapshot
		close(awaitProducer)
		delete(j.waiters, awaitConsumer)
	}
}

// shutdown close all downstream, waiters, and channels
func (j *Jet) shutdown() {
	for consumer, producer := range j.downstream {
		close(producer)
		delete(j.downstream, consumer)
	}
	for awaitConsumer, awaitProducer := range j.waiters {
		awaitProducer <- j.latestSnapshot
		close(awaitProducer)
		delete(j.waiters, awaitConsumer)
	}
	close(j.awaiter)
	close(j.registrar)
	close(j.unregistrar)
	close(j.acid)
	close(j.upstream)
}

// Up pushes a new value into the Jet
func (j *Jet) Up(data Any) {
	if j.isDone {
		return
	}

	j.upstream <- data
}

// Close shutdown the entire Jet and all downstream from Sink
func (j *Jet) Close() {
	if j.isDone {
		return
	}

	go late(func() {
		j.acid <- Signal{}
	})
}

// Sink registers a consumer channel and return it
func (j *Jet) Sink() streaming.Consumer {
	consumer := make(chan Any)

	if j.isDone {
		defer close(consumer)
	} else {
		j.registrar <- consumer
	}

	return consumer
}

// Detach unregisters a consumer channel and return an error
func (j *Jet) Detach(ch streaming.Consumer) error {
	if j.isDone {
		return errors.New("jet 'Unlink': Jet has finished or been shutdown forcefully")
	}
	j.unregistrar <- ch
	return nil
}

// Snapshots register a consumer channel and unregister on finished context
func (j *Jet) Snapshots(ctx context.Context) <-chan Any {
	sink := j.Sink()
	go func() {
		<-ctx.Done()
		_ = j.Detach(sink)
	}()
	return sink
}

// OnSnapshot register a sink and iterator over it with a callback until the provided context finishes
func (j *Jet) OnSnapshot(ctx context.Context, callback func(snapshot Any)) <-chan Signal {
	ch := j.Snapshots(ctx)
	done := make(chan Signal)

	go func() {
		for snapshot := range ch {
			callback(snapshot)
		}
		done <- Signal{}
	}()

	return done
}

// On register sink and iterate over it and call the callback
func (j *Jet) On(callback func(Any)) (<-chan Signal, func()) {
	ch := j.Sink()
	done := make(chan Signal)

	go func() {
		for snapshot := range ch {
			callback(snapshot)
		}
		done <- Signal{}
	}()

	return done, func() {
		if err := j.Detach(ch); err != nil {
			log.Println(err.Error())
		}
	}
}

// Await is method for waiting for the next value in the Jet otherwise use the latestSnapshot
func (j *Jet) Await() Any {
	res := j.AwaitNoCache()
	if res == nil {
		return j.latestSnapshot
	}
	return res
}

// AwaitNoCache is a method for waiting for the next value in the Jet but doesn't use the latestSnapshot
func (j *Jet) AwaitNoCache() Any {
	if j.isDone {
		return nil
	}

	await := make(chan Any)
	j.awaiter <- await
	return <-await
}

func (j *Jet) Done() <-chan Signal {
	done := make(chan Signal)
	go func() {
		for !j.isDone {
		}
		done <- Signal{}
	}()
	return done
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
func (j *Jet) Value() Any {
	return j.latestSnapshot
}

// Err return the accumulated error from the Jet iterator
func (j *Jet) Err() error {
	return j.accumulatedError
}
