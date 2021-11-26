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
    s := gf.NewRPCServer()
    s.RegisterFunc("mul", mul)

    go func() {
        var (
            r, w int
        )
        r = 12345
        cli := gf.NewRPCClient()
        cli.Connect("tcp", "localhost:7788")
        for {
            log.Println("发起请求........................................")
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
