//
//  types.go
//  types
//
//  Created by d-exclaimation on 6:05 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package types

// Any is just anything
type Any interface{}

// Function is any function that return a value and error (following go's standard)
type Function func() (Any, error)

// Signal is a data structure for simulating notification with little memory allocation
type Signal struct{}

