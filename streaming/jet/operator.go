//
//  operator.go
//  jet
//
//  Created by d-exclaimation on 1:58 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package jet

import . "github.com/d-exclaimation/gocurrent/types"

// Map is an operator for mapping the inner streaming value of the Jet
func Map(jt *Jet, mapper func(Any) Any) *Jet {
	newJet := New()

	// Wait for finish signal from the new Jet
	go func() {
		<-newJet.Done()
		jt.Close()
	}()

	// Iterate over the current jet and close once done
	go func() {
		ch := jt.Sink()
		for snapshot := range ch {
			newJet.Up(mapper(snapshot))
		}
		_ = jt.Detach(ch)
		newJet.Close()
	}()

	return newJet
}

// Filter is an operator for filtering the inner streaming value of the Jet
func Filter(jt *Jet, predicate func(Any) bool) *Jet {
	newJet := New()

	// Wait for finish signal from the new Jet
	go func() {
		<-newJet.Done()
		jt.Close()
	}()

	// Iterate over the current jet and close once done
	go func() {
		ch := jt.Sink()
		for snapshot := range ch {
			if predicate(snapshot) {
				newJet.Up(snapshot)
			}
		}
		_ = jt.Detach(ch)
		newJet.Close()
	}()

	return newJet
}

// FilterMap is an operator for filtering the inner streaming value of the Jet
func FilterMap(jt *Jet, predicateMap func(Any) (bool, Any)) *Jet {
	newJet := New()

	// Wait for finish signal from the new Jet
	go func() {
		<-newJet.Done()
		jt.Close()
	}()

	// Iterate over the current jet and close once done
	go func() {
		ch := jt.Sink()
		for snapshot := range ch {
			ok, res := predicateMap(snapshot)
			if ok {
				newJet.Up(res)
			}
		}
		_ = jt.Detach(ch)
		newJet.Close()
	}()

	return newJet
}
