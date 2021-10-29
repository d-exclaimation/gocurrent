//
//  lazy.go
//  jet
//
//  Created by d-exclaimation on 5:57 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package jet

import (
	"github.com/d-exclaimation/gocurrent/future"
	"github.com/d-exclaimation/gocurrent/streaming"
	. "github.com/d-exclaimation/gocurrent/types"
)

type RunnableJet func() *Jet

// Lazy setups a function to run a Jet stream
func Lazy(opts ...Option) RunnableJet {
	var (
		upstream   = make(chan Any, 2)
		register   = make(chan chan Any)
		unregister = make(chan streaming.Consumer)
		acid       = make(chan Signal)
	)

	// Setup for optional fields and configuration
	for _, opt := range opts {
		switch opt.(type) {
		case bufferedAll:
			buffer := opt.(bufferedAll)
			upstream = make(chan Any, buffer)
			register = make(chan chan Any, buffer)
			unregister = make(chan streaming.Consumer, buffer)
			acid = make(chan Signal, buffer)
		case upstreamBuffered:
			buffer := opt.(upstreamBuffered)
			upstream = make(chan Any, buffer)
			acid = make(chan Signal, buffer)
		case downstreamBuffered:
			buffer := opt.(downstreamBuffered)
			register = make(chan chan Any, buffer)
			unregister = make(chan streaming.Consumer, buffer)
			acid = make(chan Signal, buffer)
		default:
			upstream = make(chan Any)
			register = make(chan chan Any)
			unregister = make(chan streaming.Consumer)
			acid = make(chan Signal)
		}
	}

	jt := &Jet{
		_upstream:   upstream,
		_register:   register,
		_unregister: unregister,
		_await:      make(chan chan Any),
		_acid:       acid,
		downstream:  make(streaming.Downstreams),
		waiters:     make(streaming.Downstreams),
	}
	return func() *Jet {
		jt.behavior()
		return jt
	}
}

// LazyFuture setups a function to run a Jet stream with a value after future completed and closes
func LazyFuture(fut func() *future.Future) RunnableJet {
	run := Lazy()
	return func() *Jet {
		jt := run()
		go func() {
			data, err := fut().Await()
			if err == nil {
				jt.Up(data)
			}
			jt.Close()
		}()
		return jt
	}
}
