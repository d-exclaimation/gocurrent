//
//  state.go
//  pipe
//
//  Created by d-exclaimation on 10:09 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package pipe

import (
	"github.com/d-exclaimation/gocurrent/streaming/jet"
	"github.com/d-exclaimation/gocurrent/task"
	. "github.com/d-exclaimation/gocurrent/types"
)

func Seq(jt *jet.Jet) *task.Task[Any] {
	ch := jt.Sink()
	return task.Async[Any](func() (Any, error) {
		var seq []Any
		for snapshot := range ch {
			seq = append(seq, snapshot)
		}
		return seq, nil
	})
}

func Last(jt *jet.Jet) *task.Task[Any] {
	ch := jt.Sink()
	return task.Async[Any](func() (Any, error) {
		var res Any
		for snapshot := range ch {
			res = snapshot
		}
		return res, nil
	})
}

func Reduce(jt *jet.Jet, reducer func(Any, Any) Any) *task.Task[Any] {
	ch := jt.Sink()
	return task.Async[Any](func() (Any, error) {
		var res Any = nil
		for snapshot := range ch {
			if res == nil {
				res = snapshot
			} else {
				res = reducer(res, snapshot)
			}
		}
		return res, nil
	})
}
