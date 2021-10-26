//
// main.go
// futures-n-streams
//
// Created by d-exclaimation on 00:00.
//

package main

import (
	"github.com/d-exclaimation/gocurrent/future"
	"github.com/d-exclaimation/gocurrent/result"
	"log"
	"time"
)

func main() {
	fut0 := future.Async(func() (future.Any, error) {
		time.Sleep(2 * time.Second)
		return 10, nil
	})

	fut1 := future.Map(fut0, func(data future.Any) future.Any {
		intVal := data.(int)
		return intVal + 1
	})

	fut2 := future.FlatMap(fut1, func(data future.Any) *future.Future {
		intVal := data.(int)
		return slowDouble(intVal)
	})

	<-future.OnComplete(fut2, result.Case{
		Success: func(data interface{}) { log.Printf("%v \n", data) },
		Failure: func(err error) { log.Fatalln(err) },
	})

	ch := <-fut2.Channel()
	res, err := ch.Get()
	log.Printf("%v (%v)\n", res, err)

	prom := future.Maybe()
	go func() {
		time.Sleep(1 * time.Second)
		prom.Success(10)
	}()

	res1 := result.Await(prom.Future())

	res2 := result.Map(res1, func(i interface{}) interface{} { return i.(int) + 1 })

	res3 := result.Filter(res2, func(i interface{}) bool { return i.(int) < 10 })

	res3.Match(result.Case{
		Success: func(i interface{}) { log.Println(i) },
		Failure: func(err error) { log.Fatalln(err) },
	})
}

func slowDouble(intVal int) *future.Future {
	return future.Async(func() (future.Any, error) {
		time.Sleep(2 * time.Second)
		return intVal + 10, nil
	})
}
