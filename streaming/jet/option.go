//
//  option.go
//  jet
//
//  Created by d-exclaimation on 10:11 PM.
//  Copyright © 2021 d-exclaimation. All rights reserved.
//

package jet

import "github.com/d-exclaimation/gocurrent/types"

// Option is a interface pattern to be used for constructing Jet streams
type Option interface {
	// implement is a required method for allowing any settings to follow Option
	implement()
}

// upstreamBuffered is an Option for Jet stream with specified buffer size on upstream
type upstreamBuffered int

func (u upstreamBuffered) implement() {}

// WithUpstreamBuffer is an Option to add buffer to Jet stream's upstream
func WithUpstreamBuffer(buffer int) Option {
	return upstreamBuffered(buffer)
}

// bufferedAll is an Option for Jet stream with specified buffer size for all channel
type bufferedAll int

func (b bufferedAll) implement() {}

// WithBuffer is an Option to add buffer to all Jet stream's channel
func WithBuffer(buffer int) Option {
	return bufferedAll(buffer)
}

// downstreamBuffered is an Option for Jet stream with specified buffer size for downstream
type downstreamBuffered int

func (d downstreamBuffered) implement() {}

// WithDownstreamBuffer is an Option to add buffer to all Jet stream's downstream
func WithDownstreamBuffer(buffer int) Option {
	return downstreamBuffered(buffer)
}

// notBuffered is an Option for Jet stream with no buffer
type notBuffered types.Signal

func (d notBuffered) implement() {}

// WithNoBuffer is an Option to remove buffer to all Jet stream channel
func WithNoBuffer() Option {
	return notBuffered{}
}
