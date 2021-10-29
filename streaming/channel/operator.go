//
//  operator.go
//  consumer
//
//  Created by d-exclaimation on 2:57 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package channel

import (
	"context"
	"github.com/d-exclaimation/gocurrent/streaming"
	. "github.com/d-exclaimation/gocurrent/types"
)

// Map adds a pipeline function on to the channel result
func Map(ch streaming.Consumer, mapper func(interface{}) interface{}) streaming.Consumer {
	channel := make(chan Any)
	go func() {
		for incoming := range ch {
			channel <- mapper(incoming)
		}
	}()
	return channel
}

// ApplyContext applies all the necessary setup with the context for closing and receiving data in channel
func ApplyContext(ch streaming.Consumer, ctx context.Context) streaming.Consumer {
	outgoing := make(chan Any)
	bridge := make(chan interface{})
	acid := make(chan struct{})

	go func() {
		<-ctx.Done()
		acid <- struct{}{}
	}()

	go func() {
		for {
			select {
			case incoming := <-bridge:
				outgoing <- incoming
			case _ = <-acid:
				close(outgoing)
				return
			}
		}
	}()

	go func() {
		for i := range ch {
			bridge <- i
		}
		acid <- struct{}{}
	}()

	return outgoing
}
