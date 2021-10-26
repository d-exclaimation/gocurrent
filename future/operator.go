//
//  operator.go
//  future
//
//  Created by d-exclaimation on 1:33 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package future

import (
	"errors"
	"github.com/d-exclaimation/gocurrent/result"
)

// Map is an operator for mapping the inner successful value of the Future
func Map(f *Future, mapper func(Any) Any) *Future {
	return Async(func() (Any, error) {
		data, err := f.Await()
		if err != nil {
			return nil, err
		}
		return mapper(data), err
	})
}

// FlatMap is another operator for mapping the inner value to another Future and flatten it
func FlatMap(f *Future, mapper func(Any) *Future) *Future {
	return Async(func() (Any, error) {
		data, err := f.Await()
		if err != nil {
			return nil, err
		}
		return mapper(data).Await()
	})
}

// OnComplete is another operator for performing callback on the Future result
// and return a channel to wait for that finished callback
func OnComplete(f *Future, callback result.Case) <-chan bool {
	await := make(chan bool)
	go func() {
		data, err := f.Await()
		res := result.New(data, err)
		res.Match(callback)
		await <- true
		close(await)
	}()
	return await
}

// Filter is an operator for filter the inner value
func Filter(r *Future, predicate func(interface{}) bool) *Future {
	return Async(func() (Any, error) {
		data, err := r.Get()
		if err != nil {
			return nil, err
		}
		if !predicate(data) {
			return nil, errors.New("future 'Filter': Value doesn't met the predicate")
		}
		return data, nil
	})
}
