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
	"github.com/d-exclaimation/gocurrent/streaming/common"
	. "github.com/d-exclaimation/gocurrent/types"
)

type RunnableJet func() *Jet

// Lazy setups a function to run a Jet stream
func Lazy() RunnableJet {
	jt := &Jet{
		_upstream:        make(chan Any),
		_register:        make(chan chan Any),
		_unregister:      make(chan common.Consumer),
		_await:           make(chan chan Any),
		_acid:            make(chan Signal),
		latestSnapshot:   nil,
		accumulatedError: nil,
		downstream:       make(common.Downstreams),
		waiters:          make(common.Downstreams),
		isDone:           false,
	}
	return func() *Jet {
		jt.behavior()
		return jt
	}
}

// LazyJust setups a function to run a Jet stream with a single value and closes
func LazyJust(single func() interface{}) RunnableJet {
	run := Lazy()
	return func() *Jet {
		jt := run()
		jt.Up(single())
		jt.Close()
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
