package main

import (
    "context"
    gf "github.com/DGHeroin/gofaster"
    "log"
)

func main() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    s := gf.NewRPCServer()
    s.RegisterFunc("mul", mul)
    s.StartServe("tcp", ":7788")
}
func mul(ctx context.Context, r *int, w *int) error {
//    log.Println("收到", *r)

    *w = *r + 10
    //time.Sleep(time.Second*5)
    return nil //fmt.Errorf("处理错误")
}
