package main

import (
    "context"
    fast "github.com/DGHeroin/fast"
    "go.uber.org/ratelimit"
    "log"
    "sync/atomic"
    "time"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    var last int64
    go func() {
        for {
            time.Sleep(time.Second)
            now := atomic.LoadInt64(&count)
            qps := now - last
            last = now
            log.Println("==qps:", qps)
        }
    }()
    for i := 0; i < 50; i++ {
        go doTest()
    }
    select {}
}
func doTest() {
    var (
        r, w int
    )
    r = 12345
    cli := fast.NewRPCClient()

    cli.OnEvent(func(event fast.RPCEvent, args ...interface{}) {
        switch event {
        case fast.EventOpen:
            cli.AddPlugin(ratelimit.New(1000))
        }
    })
    cli.Connect("tcp", "localhost:7788")

    n := 2
    for {
        n--
        if !cli.IsConnected() {
            time.Sleep(time.Second)
            continue
        }

        cli.Call(context.Background(), "mul", &r, &w)
        atomic.AddInt64(&count, 1)

        //cli.Go(context.Background(), "mul", &r, &w, func(err error) {
        //   atomic.AddInt64(&count, 1)
        //})
    }
}

var count int64
