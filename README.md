# gocurrent

## Future & Promise

Deferred value that may be present in the future.

```go
package main

import (
    "github.com/d-exclaimation/gocurrent/future"
    "github.com/d-exclaimation/gocurrent/result"
    "log"
    "time"
)

func main() {
    fut0 := future.Async(func() (future.Any, error) {
        time.Sleep(10 * time.Second)
        return 10, nil
    })

    fut1 := future.Map(fut0, func(any future.Any) future.Any {
        return any.(int) + 2
    })

    future.OnComplete(fut1, result.Case{
        Success: func(i interface{}) {
            log.Println(i)
        },
        Failure: func(err error) {
            log.Fatalln(err)
        },
    })

    // ... wait for 10 seconds
    // log: 12
}
```

## Jet streams

Time based single topic data stream with a single upstream and multiple downstream (Broadcast / hot stream).

```go
package main

import (
    "github.com/d-exclaimation/gocurrent/streaming/jet"
    "log"
)

func main() {
    jt := jet.New()

    jt.OnSnapshot(func(i interface{}) {
        log.Printf("[1]: %v, ", i)
    })

    go func() {
        ch := jt.Sink()
        for {
            select {
            case i := <-ch:
                log.Printf("[2]: %v, ", i)
            }
        }
    }()

    go func() {
        for jt.Next() {
            log.Printf("[3]: %v\n", jt.Value())
        }
    }()

    jt.Up(1)
    // log: [1]: 1, [2]: 1, [3]: 1
    jt.Up(2)
    // log: [1]: 2, [2]: 2, [3]: 2
    jt.Up(3)
    // log: [1]: 2, [2]: 2, [3]: 2
    jt.Close()
}
```
