package main

import (
    "flag"
    gof "github.com/DGHeroin/gofaster"
    "log"
    "net"
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
    network = flag.String("n", "kcp", "network type")
)

func (c *clientHandler) OnOpen(client *gof.Client) {
    log.Println("OnOpen", client)
    n := 0
    for n < 5 {
        sendData := []byte(time.Now().Format(time.RFC3339))
        log.Println("client send", sendData)
        client.SendPacket(sendData)
        time.Sleep(time.Second)
        n++
    }
    client.Close()
}


func (c *clientHandler) HandlePacket(client *gof.Client,packet *gof.Packet) {
    log.Println("Client============>HandlePacket", packet.Payload())
}

func main() {
    flag.Parse()
    log.SetFlags(log.LstdFlags|log.Lshortfile)
    startServer()
}

func startClient() {
    time.Sleep(time.Second)
    handler := &clientHandler{}
    client := gof.NewClient(*network, handler, gof.Option{Address: "127.0.0.1:5566"})
    defer client.Close()
    time.Sleep(time.Second * 10)
}

func startServer() {
    handler := &serverHandler{}
    gof.Serve(*network, handler, gof.Option{Address: ":5566"})
}

func (s *serverHandler) OnStartServe(addr net.Addr) {
    log.Println("start service", addr.String())
    gof.GoN(1, startClient)
}

func (s *serverHandler) HandlePacket(client *gof.Client, packet *gof.Packet) {
    log.Println("Server------------>HandlePacket", packet.Payload())
    client.SendPacket([]byte("hello world"))
}
