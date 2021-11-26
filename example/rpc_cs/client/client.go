package main

import (
    "context"
    gf "github.com/DGHeroin/gofaster"
    "log"
    "sync/atomic"
    "time"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    var (
        r, w int
    )
    r = 12345
    cli := gf.NewRPCClient()
    cli.Connect("tcp", "localhost:7788")

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

    for {
        if !cli.IsConnected() {
            time.Sleep(time.Second)
            log.Println("等..")
            continue
        }
        startTime := time.Now()
        err := cli.Call(context.Background(), "mul", &r, &w)
        if err != nil {

        }
        log.Println("==>", time.Since(startTime))
        break
        //if err != nil {
        //    log.Println("搞错了", err)
        //} else {
        //    log.Println("结果是", w)
        //}
        atomic.AddInt64(&count, 1)
        //time.Sleep(time.Second)
    }
}
var count int64