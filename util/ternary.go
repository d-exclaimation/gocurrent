//
//  ternary.go
//  util
//
//  Created by d-exclaimation on 7:07 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package util

// Ternary is a ternary operator in function form. A bit more complicated; however, it functions nonetheless
func Ternary[T any](predicate bool, truthy func() T, falsy func() T) T {
	if predicate {
		return truthy()
	}
	return falsy()
}
