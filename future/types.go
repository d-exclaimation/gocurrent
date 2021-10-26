//
//  types.go
//  future
//
//  Created by d-exclaimation on 3:07 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package future

// DeliveryStatus indicate the progress of a certain function
type DeliveryStatus string

const (
	// Idle status show function hasn't executed
	Idle DeliveryStatus = "future-idle"

	// Loading status show function has executed but hasn't finished
	Loading DeliveryStatus = "future-loading"

	// Success status show function has finished executing and return a successful value
	Success DeliveryStatus = "future-success"

	// Failure status show function has finished executing and return an unsuccessful value
	Failure DeliveryStatus = "future-failure"
)

// Any is just anything
type Any interface{}

// Function is any function that return a value and error (following go's standard)
type Function func() (Any, error)
