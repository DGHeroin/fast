package main

import (
    gof "github.com/DGHeroin/gofaster"
    "log"
    "net"
    "strconv"
    "sync/atomic"
    "time"
)

type (
    serverHandler struct{}
)

var (
    count int32
)

func main() {
    startServer()
}

func startClient() {
    time.Sleep(time.Second)
    client, err := gof.NewClient("tcp", gof.Option{Address: "127.0.0.1:5566"})
    if err != nil {
        log.Println(err)
        return
    }
    defer client.Close()
    for {
        client.SendPacket([]byte(time.Now().Format(time.RFC3339)))
    }

}

func startServer() {
    go func() {
        lastCount := count
        for {
            nowCount := atomic.LoadInt32(&count)
            qps := nowCount - lastCount
            lastCount = nowCount
            log.Println("qps", Format3(int64(qps)))
            time.Sleep(time.Second)
        }
    }()
    handler := &serverHandler{}
    gof.Serve("tcp", handler, gof.Option{Address: ":5566"})
}

func (s *serverHandler) OnStartServe(addr net.Addr) {
    log.Println("start service", addr.String())
    gof.GoN(10, startClient)
}

func (s *serverHandler) HandlePacket(packet *gof.Packet) {
    atomic.AddInt32(&count, 1)
}
func (s *serverHandler) OnNew(conn net.Conn) {
    log.Println("新链接", conn.RemoteAddr())
}

func (s *serverHandler) OnClose(conn net.Conn) {
    log.Println("关闭链接", conn.RemoteAddr())
}
func Format3(n int64) string {
    if n < 0 {
        return "-" + Format3(-n)
    }
    in := []byte(strconv.FormatInt(n, 10))

    var out []byte
    if i := len(in) % 3; i != 0 {
        if out, in = append(out, in[:i]...), in[i:]; len(in) > 0 {
            out = append(out, ',')
        }
    }
    for len(in) > 0 {
        if out, in = append(out, in[:3]...), in[3:]; len(in) > 0 {
            out = append(out, ',')
        }
    }
    return string(out)
}
