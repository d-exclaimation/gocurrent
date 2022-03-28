//
//  instance.go
//  jet
//
//  Created by d-exclaimation on 4:27 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package jet

import (
	"github.com/d-exclaimation/gocurrent/task"
	. "github.com/d-exclaimation/gocurrent/types"
	"time"
)

// From instantiate a new Jet stream from a channel
func From(ch <-chan Any, opts ...Option) *Jet {
	jt := New(opts...)
	bridge := make(chan Any)
	acid := make(chan struct{})

	// End signal from channel
	go func() {
		for incoming := range ch {
			bridge <- incoming
		}
		acid <- struct{}{}
	}()

	// End signal from Jet
	go func() {
		<-jt.Done()
		acid <- struct{}{}
	}()

	// Main iteration go routine
	go func() {
		for {
			select {
			case incoming := <-bridge:
				jt.Up(incoming)
			case _ = <-acid:
				jt.Close()
				return
			}
		}
	}()

	return jt
}

// Empty instantiate a Jet stream that push no value and closes
func Empty() *Jet {
	jt := New()
	defer jt.Close()
	return jt
}

// Future instantiate a Jet stream with a value after future completed and closes
func Future(fut *task.Task[Any], opts ...Option) *Jet {
	jt := New(opts...)
	defer routine(func() {
		data, err := fut.Await()
		if err == nil {
			jt.Up(data)
		}
		defer jt.Close()
	})
	return jt
}

// routine allocate function late-ly to a new goroutine
func routine(fn func()) {
	go late(fn)
}

// late call a function lately
func late(fn func()) {
	time.Sleep(10 * time.Millisecond)
	fn()
}
