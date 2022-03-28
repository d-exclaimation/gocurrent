//
//  operator.go
//  task
//
//  Created by d-exclaimation on 10:03 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package task

import "errors"

// Map transformed a wrapped value of a Task into a new type
func Map[T, K any](t *Task[T], transform func(T) (K, error), base K) *Task[K] {
	return Async[K](func() (K, error) {
		res, err := t.Await()
		if err != nil {
			return base, err
		}
		return transform(res)
	})
}

// FlatMap transformed a wrapped value of a Task into a Task with new type
func FlatMap[T, K any](t *Task[T], transform func(T) *Task[K], base K) *Task[K] {
	return Async[K](func() (K, error) {
		res, err := t.Await()
		if err != nil {
			return base, nil
		}
		return transform(res).Await()
	})
}

// Filter filters the wrapped value or throw an error
func Filter[T any](r *Task[T], predicate func(T) bool) *Task[T] {
	return Async[T](func() (T, error) {
		data, err := r.Await()
		if err != nil {
			return data, err
		}
		if !predicate(data) {
			return data, errors.New("future 'Filter': Value doesn't met the predicate")
		}
		return data, nil
	})
}

// Recover recovers an error into the wrapped value or throw an error
func Recover[T any](r *Task[T], recovery func(err error) (T, error)) *Task[T] {
	return Async[T](func() (T, error) {
		data, err := r.Await()
		if err != nil {
			return recovery(err)
		}
		return data, nil
	})
}
