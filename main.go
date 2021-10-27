//
// main.go
// futures-n-streams
//
// Created by d-exclaimation on 00:00.
//

package main

import (
	"github.com/d-exclaimation/gocurrent/future"
	"github.com/d-exclaimation/gocurrent/streaming/jet"
	. "github.com/d-exclaimation/gocurrent/types"
	"log"
	"time"
)

func main() {
	fut0 := future.Async(func() (Any, error) {
		time.Sleep(time.Second * 5)
		return 1, nil
	})

	fut1 := future.Map(fut0, func(any Any) Any {
		return any.(int) + 1
	})

	jt := jet.Future(fut1)

	<-jt.Done()
	log.Println(jt.Await())
}
