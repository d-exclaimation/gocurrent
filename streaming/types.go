//
//  types.go
//  jet
//
//  Created by d-exclaimation on 9:44 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package streaming

import . "github.com/d-exclaimation/gocurrent/types"

// Consumer is a channel that only allow consuming the data
type Consumer <-chan Any

// Producer is a channel that only allow pushing data
type Producer chan<- Any

// Downstreams is a store for consumer-producer pair
type Downstreams map[Consumer]Producer
