//
//  try.go
//  try
//
//  Created by d-exclaimation on 6:41 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package try

// Try is a data structure representing a data, error pair
type Try[T any] struct {
	value T
	error error
}

// From instantiate a new Try from throwing function
func From[T any](op func() (T, error)) *Try[T] {
	value, err := op()
	return &Try[T]{value, err}
}

// New instantiate a new Try from a tuple
func New[T any](value T, err error) *Try[T] {
	return &Try[T]{value, err}
}

// Get return the successful value regardless of failures
func (try *Try[T]) Get() T {
	return try.value
}

// ToOption return the inner values as tuples
func (try *Try[T]) ToOption() (T, error) {
	return try.value, try.error
}

// Match tries to match the wrapped value to the proper cases
func (try *Try[T]) Match(matcher Case[T]) {
	if try.IsSuccess() {
		matcher.Success(try.Get())
	} else {
		matcher.Failure(try.error)
	}
}

// IsSuccess indicates whether the Try succeed
func (try *Try[T]) IsSuccess() bool {
	return try.error == nil
}

// OrElse return the successful value if succeeded, or the fallback
func (try *Try[T]) OrElse(fallback T) T {
	if try.IsSuccess() {
		return try.Get()
	}
	return fallback
}
