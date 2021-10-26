//
// main.go
// futures-n-streams
//
// Created by d-exclaimation on 00:00.
//

package main

import (
	"github.com/d-exclaimation/gocurrent/streaming/jet"
	"log"
	"time"
)

func main() {
	jt := jet.New()

	doneAll := make(chan bool)
	go func() {
		done1, _ := jt.OnSnapshot(func(i interface{}) {
			log.Printf("[1]: %v\n", i)
		})

		done2, _ := jt.OnSnapshot(func(i interface{}) {
			log.Printf("[2]: %v\n", i)
		})

		f1 := <-done1
		f2 := <-done2

		doneAll <- f1 && f2
	}()

	go func() {
		for jt.Next() {
			log.Printf("[3]: %v\n", jt.Value())
		}
	}()

	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		jt.Up(i)
	}
	jt.Close()

	<-doneAll
}
