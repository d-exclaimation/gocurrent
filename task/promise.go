//
//  promise.go
//  task
//
//  Created by d-exclaimation on 8:24 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package task

import (
	"github.com/d-exclaimation/gocurrent/try"
)

type Promise[T any] struct {
	job     *Task[T]
	mailbox chan *try.Try[T]
}

// Maybe construct a new Promise
func Maybe[T any]() *Promise[T] {
	delivery := make(chan *try.Try[T])
	prom := &Promise[T]{
		job: Async[T](func() (T, error) {
			res := <-delivery
			return res.ToOption()
		}),
		mailbox: delivery,
	}
	return prom
}

// Task return the inner task but doesn't allow mutation
func (p *Promise[T]) Task() *Task[T] {
	return p.job
}

// Success finishes the future with a successful value
func (p *Promise[T]) Success(data T) {
	p.mailbox <- try.New(data, nil)
}

// Failure finishes the future with an unsuccessful value (base is used for non-nullable T).
func (p *Promise[T]) Failure(err error, base T) {
	p.mailbox <- try.New(base, err)
}
