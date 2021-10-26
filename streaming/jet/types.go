//
//  types.go
//  jet
//
//  Created by d-exclaimation on 9:44 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package jet

type Consumer <-chan interface{}

type Producer chan<- interface{}

type Downstreams map[Consumer]Producer

type OneTimers map[chan interface{}]bool
