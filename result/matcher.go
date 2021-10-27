//
//  matcher.go
//  result
//
//  Created by d-exclaimation on 2:00 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package result

import . "github.com/d-exclaimation/gocurrent/types"

// Case is a object for performing Result pattern matching
type Case struct {
	// Success case
	Success func(Any)

	// Failure case
	Failure func(error)
}
