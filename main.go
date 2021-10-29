//
// main.go
// futures-n-streams
//
// Created by d-exclaimation on 00:00.
//

package main

import (
	"github.com/d-exclaimation/gocurrent/streaming/jet"
	"github.com/d-exclaimation/gocurrent/types"
	"log"
)

func main() {
	jt := jet.New()

	jt.On(func(any types.Any) {
		log.Println(any)
	})

	for i := 0; i < 10; i++ {
		jt.Up(i)
	}
	jt.Close()

	<-jt.Done()
}
