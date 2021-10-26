//
//  promise.go
//  promise
//
//  Created by d-exclaimation on 3:32 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package future

import (
	"github.com/d-exclaimation/gocurrent/result"
)

// Promise is a data structure to create dynamically finished Future
type Promise struct {
	// fut is the inner Future for this Promise
	fut *Future

	// _Delivery is the channel to suspend the Future callback push the resulting value
	_Delivery chan result.Result

	// isClosed indicate where _Delivery has been closed
	isClosed bool
}

// Maybe construct a new Promise
func Maybe() *Promise {
	delivery := make(chan result.Result)
	prom := &Promise{
		fut: Async(func() (Any, error) {
			res := <-delivery
			return res.Get()
		}),
		_Delivery: delivery,
		isClosed:  false,
	}
	return prom
}

// Future return the inner future but doesn't allow mutation
func (p *Promise) Future() *Future {
	return p.fut
}

// Success finishes the future with a successful value
func (p *Promise) Success(data Any) {
	if p.isClosed {
		return
	}
	p._Delivery <- result.New(data, nil)
	p.close()
}

// Failure finishes the future with an unsuccessful value.
func (p *Promise) Failure(err error) {
	if p.isClosed {
		return
	}
	p._Delivery <- result.New(nil, err)
	p.close()
}

// close the _Delivery channel and set the isClosed properly
func (p *Promise) close() {
	close(p._Delivery)
	p.isClosed = true
}
