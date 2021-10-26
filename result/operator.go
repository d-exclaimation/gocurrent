//
//  operator.go
//  result
//
//  Created by d-exclaimation on 4:16 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package result

import "errors"

// Map is an operator for mapping the inner successful value of the Result
func Map(r Result, mapper func(interface{}) interface{}) Result {
	data, err := r.Get()
	if err != nil {
		return New(nil, err)
	}
	return New(mapper(data), nil)
}

// FlatMap is an operator for mapping the inner value to another Result
func FlatMap(r Result, mapper func(interface{}) Result) Result {
	data, err := r.Get()
	if err != nil {
		return New(nil, err)
	}
	return mapper(data)
}

// ForEach is an operator for call a callback on the inner value
func ForEach(r Result, callback func(interface{})) {
	data, err := r.Get()
	if err != nil {
		return
	}
	callback(data)
}

// Filter is an operator for filter the inner value
func Filter(r Result, predicate func(interface{}) bool) Result {
	data, err := r.Get()
	if err != nil {
		return New(nil, err)
	}
	if !predicate(data) {
		return New(nil, errors.New("result 'Filter': Value doesn't met the predicate"))
	}
	return New(data, nil)
}
