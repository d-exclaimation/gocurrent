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
	"github.com/d-exclaimation/gocurrent/streaming/common"
	. "github.com/d-exclaimation/gocurrent/types"
)

func Map(ch common.Consumer, mapper func(interface{}) interface{}) common.Consumer {
	channel := make(chan Any)
	go func() {
		for incoming := range ch {
			channel <- mapper(incoming)
		}
	}()
	return channel
}

func Solicit(ch common.Consumer, ctx context.Context) common.Consumer {
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