//
//  instance.go
//  jet
//
//  Created by d-exclaimation on 4:27 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package jet

import (
	"github.com/d-exclaimation/gocurrent/future"
	. "github.com/d-exclaimation/gocurrent/types"
	"time"
)

// From instantiate a new Jet stream from a channel
func From(ch <-chan Any) *Jet {
	jt := New()
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

// Just instantiate a Jet stream with only a single value then closes
func Just(single Any) *Jet {
	jt := New()
	go func() {
		jt.Up(single)
		defer jt.Close()
	}()
	return jt
}

// Of instantiate a Jet stream with all the provided values then closes
func Of(multi ...Any) *Jet {
	jt := New()
	go func() {
		for _, each := range multi {
			jt.Up(each)
		}
		defer jt.Close()
	}()
	return jt
}

// Empty instantiate a Jet stream that push no value and closes
func Empty() *Jet {
	jt := New()
	jt.Close()
	return jt
}

// Future instantiate a Jet stream with a value after future completed and closes
func Future(fut *future.Future) *Jet {
	jt := New()
	go func() {
		data, err := fut.Await()
		if err == nil {
			jt.Up(data)
		}
		time.Sleep(time.Millisecond * 200)
		jt.Close()
	}()
	return jt
}
