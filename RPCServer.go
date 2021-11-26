package gf

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"
    "time"
)

type (
    RPCServer struct {
        ServerHandler
        mu         sync.Mutex
        registered map[string]func(ctx context.Context, req []byte) (resp []byte, err error)
    }

    RPCContext struct {
        context.Context
        values map[string]interface{}
    }
)

var (
    ErrMessageFormat   = fmt.Errorf("message format error")
    ErrHandlerNotFound = fmt.Errorf("message handler not found")
)

func NewRPCServer() *RPCServer {
    r := &RPCServer{
        registered: map[string]func(ctx context.Context, req []byte) (resp []byte, err error){},
    }

    return r
}
func (s *RPCServer) RegisterFunc(name string, fn interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    cb, err := MakeRPCFunc(name, fn, 0, 0)
    if err != nil {
        return
    }
    s.registered[name] = cb
}

func (s *RPCServer) onMessage(client *Client, msg *message) {
    switch msg.Type {
    case mTypeRequest:
        s.handleRequest(client, msg)
    }
}
func replyError(client *Client, msgId int32, err error) {
    respMsg := newMessage(mTypeResponse)
    defer respMsg.Release()
    respMsg.Id = msgId
    respMsg.Err = []byte(err.Error())
    sendData, _ := MSGPack(respMsg)

    client.SendPacket(sendData)
}

func (s *RPCServer) handleRequest(client *Client, msg *message) {
    if msg.Name == nil {
        replyError(client, msg.Id, ErrMessageFormat)
        return
    }
    s.mu.Lock()
    fn, ok := s.registered[*msg.Name]
    s.mu.Unlock()
    if !ok {
        replyError(client, msg.Id, ErrHandlerNotFound)
        return
    }
    respData, err := fn(context.Background(), msg.Payload)
    if err != nil {
        replyError(client, msg.Id, err)
        return
    }
    respMsg := newMessage(mTypeResponse)
    defer respMsg.Release()
    respMsg.Id = msg.Id
    respMsg.Payload = respData

    sendData, _ := MSGPack(respMsg)
    client.SendPacket(sendData)
}

func (s *RPCServer) StartServe(network string, address string) {
    Serve(network, s, Option{
        Address:       address,
        TLS:           nil,
        RetryDuration: time.Second,
    })
}

func (s *RPCServer) OnStartServe(addr net.Addr) {
    log.Println("启用", addr)
}
func (s *RPCServer) OnNew(client *Client) {
    log.Println("新链接..", client)
}
func (s *RPCServer) OnClose(client *Client) {
    log.Println("关闭链接..", client)
}
func (s *RPCServer) HandlePacket(client *Client, packet *Packet) {
    payload := packet.Payload()
    msg := allocMessage()
    defer msg.Release()
    err := MSGUnpack(payload, msg)
    if err != nil {
        log.Println(err)
        return
    }
    s.onMessage(client, msg)

}
