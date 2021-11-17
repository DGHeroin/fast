package main

import (
    gof "github.com/DGHeroin/gofaster"
    "log"
    "net"
    "time"
)

type (
    serverHandler struct {
    }
)

func main() {
    gof.GoN(1, startClient)
    startServer()
}

func startClient() {
    time.Sleep(time.Second)
    tlsConfig, err := gof.LoadTLSConfig("example/certs/client.pem", "example/certs/client.key", "example/certs/client.pem")
    if err != nil {
        log.Println(err)
        return
    }
    client, err := gof.NewClient("tcp", gof.Option{Address: "127.0.0.1:5566", TLS: tlsConfig})
    if err != nil {
        log.Println(err)
        return
    }
    defer client.Close()
    n := 0
    for n < 10 {
        client.SendPacket([]byte(time.Now().Format(time.RFC3339)))
        time.Sleep(time.Second)
        n++
    }

}

func startServer() {
    tlsConfig, err := gof.LoadTLSConfig("example/certs/server.pem", "example/certs/server.key", "example/certs/client.pem")
    if err != nil {
        log.Println(err)
        return
    }
    handler := &serverHandler{}
    gof.Serve("tcp", handler, gof.Option{Address: ":5566", TLS: tlsConfig})
}

func (s *serverHandler) OnStartServe(addr net.Addr) {
    log.Println("start service", addr.String())
}

func (s *serverHandler) HandlePacket(packet *gof.Packet) {
    log.Println("收到", packet.PayloadAsString())
}
func (s *serverHandler) OnNew(conn net.Conn) {
    log.Println("新链接", conn.RemoteAddr())
}

func (s *serverHandler) OnClose(conn net.Conn) {
    log.Println("关闭链接", conn.RemoteAddr())
}