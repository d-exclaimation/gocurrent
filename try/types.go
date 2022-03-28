//
//  types.go
//  task
//
//  Created by d-exclaimation on 6:50 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package try

import "log"

// Case is a object for performing Result pattern matching
type Case[T any] struct {
	// Success case
	Success func(T)

	// Failure case
	Failure func(error)
}

func Log[T any](a T) {
	log.Println(a)
}

func LogFatal(err error) {
	log.Fatalln(err.Error())
}

func LogError(err error) {
	log.Println(err.Error())
}
