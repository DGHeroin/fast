package main

import (
    "context"
    "github.com/DGHeroin/fast"
    "log"
    "sync/atomic"
    "time"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    s := fast.NewRPCServer()

    s.OnEvent(func(event fast.RPCEvent, args ...interface{}) {
        switch event {
        case fast.EventAccept:
            client := args[0].(*fast.Client)
            time.AfterFunc(time.Second, func() {
                //client.Call(context.Background(), "sayHello", nil, nil)
                client.SendPacket([]byte("hello world"))
            })
        }
    })

    s.RegisterFunc("mul", mul)

    go func() {
        var (
            r, w int
        )
        r = 12345
        cli := fast.NewRPCClient()
        cli.RegisterFunc("sayHello", func(ctx context.Context, r *int, w *int) error {
            log.Println("收到...")
            return nil
        })
        cli.OnEvent(func(event fast.RPCEvent, args ...interface{}) {
            switch event {
            case fast.EventRawMessage:
                //client := args[0].(fast.Client)
                data := args[1].([]byte)
                log.Println("==>", string(data))
            }
        })
        cli.Connect("tcp", "localhost:7788")

        for {
            if !cli.IsConnected() {
                time.Sleep(time.Millisecond)
                continue
            }
            err := cli.Call(context.Background(), "mul", &r, &w)
            if err != nil {
                log.Println("搞错了", err)
            } else {
                log.Println("结果是", w)
            }
            time.Sleep(time.Second)
        }
    }()

    s.StartServe("tcp", ":7788")

}

var count int32

func mul(ctx context.Context, r *int, w *int) error {
    log.Println("收到", *r)
    atomic.AddInt32(&count, 1)
    *w = *r + 10
    //time.Sleep(time.Second*5)
    return nil //fmt.Errorf("处理错误")
}
