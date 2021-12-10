package main

import (
    "context"
    "github.com/DGHeroin/fast"
    "go.uber.org/ratelimit"
    "log"
    "time"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    s := fast.NewRPCServer()
    s.AddPlugin(ratelimit.New(500))

    s.OnEvent(func(event fast.RPCEvent, args ...interface{}) {
        switch event {
        case fast.EventAccept:
            time.AfterFunc(time.Second, func() {

            })
        }
    })

    s.RegisterFunc("mul", mul)
    s.StartServe("tcp", ":7788")
}
func mul(ctx context.Context, r *int, w *int) error {
    //    log.Println("收到", *r)

    *w = *r + 10
    //time.Sleep(time.Second*5)
    return nil //fmt.Errorf("处理错误")
}
