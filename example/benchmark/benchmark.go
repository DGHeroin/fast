package main

import (
    "flag"
    gof "github.com/DGHeroin/fast"
    "log"
    "net"
    "strconv"
    "sync/atomic"
    "time"
)

type (
    serverHandler struct {
        gof.ServerHandler
    }
    clientHandler struct {
        gof.ClientHandler
    }
)

var (
    count   int32
    network = flag.String("n", "kcp", "network type")
)

func main() {
    flag.Parse()
    startServer()
}
func (clientHandler) OnOpen(client *gof.Client) {
    go func() {
        for {
            client.SendPacket([]byte(time.Now().Format(time.RFC3339)))
        }
    }()
}
func startClient() {
    time.Sleep(time.Second)
    handler := &clientHandler{}
    client := gof.NewClient(*network, handler, gof.Option{Address: "127.0.0.1:5566"})
    defer client.Close()
    for {
        time.Sleep(time.Second)
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
    gof.Serve(*network, handler, gof.Option{Address: ":5566"})
}

func (s *serverHandler) OnStartServe(addr net.Addr) {
    log.Println("start service", addr.String())
    gof.GoN(100, startClient)
}

func (s *serverHandler) HandlePacket(client *gof.Client, packet *gof.Packet) {
    atomic.AddInt32(&count, 1)
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
