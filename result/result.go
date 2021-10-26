//
//  result.go
//  result
//
//  Created by d-exclaimation on 1:57 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package result

import "github.com/d-exclaimation/gocurrent/future"

// Result is a data structure representing a data, error pair
type Result struct {
	// Value is the inner value of the result
	Value interface{}

	// Err is the inner error of the result
	Err error
}

// Match tries to match the wrapped value to the proper cases
func (r *Result) Match(cases Case) {
	if r.Err != nil {
		cases.Failure(r.Err)
		return
	}
	cases.Success(r.Value)
}

// Get return the inner values as tuples
func (r *Result) Get() (interface{}, error) {
	return r.Value, r.Err
}

// Recover catch the error and transform it to the successful value
func (r *Result) Recover(fallback func(error) interface{}) interface{} {
	if r.Err != nil {
		return fallback(r.Err)
	}
	return r.Value
}

// IsSuccess return a boolean indicate whether result is successful
func (r *Result) IsSuccess() bool {
	return r.Err == nil
}

// New instantiate a new Result
func New(data interface{}, err error) Result {
	return Result{
		Value: data,
		Err:   err,
	}
}

// From convert a go's standard throwable function return value into a Result
func From(run func() (interface{}, error)) Result {
	data, err := run()
	return New(data, err)
}

// Await awaits and converts a Future return value into a Result
func Await(fut *future.Future) Result {
	return From(func() (interface{}, error) {
		return fut.Await()
	})
}
