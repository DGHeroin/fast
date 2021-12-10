package fast

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
        rpcMap    RPCMap
        rateLimit RateLimit
        onEvent   func(event RPCEvent, args ...interface{})
    }
    RPCMap struct {
        mu         sync.Mutex
        registered map[string]func(ctx context.Context, req []byte) (resp []byte, err error)
    }

    RateLimit interface {
        Take() time.Time
    }

    RPCContext struct {
        context.Context
        values map[string]interface{}
    }
)

func (h *RPCMap) init() {
    h.registered = map[string]func(ctx context.Context, req []byte) (resp []byte, err error){}
}
func (s *RPCMap) RegisterFunc(name string, fn interface{}) {
    s.mu.Lock()
    defer s.mu.Unlock()
    cb, err := MakeRPCFunc(name, fn, 0, 0)
    if err != nil {
        return
    }
    s.registered[name] = cb
}
func (s *RPCMap) Get(name string) func(ctx context.Context, req []byte) (resp []byte, err error) {
    s.mu.Lock()
    fn, _ := s.registered[name]
    s.mu.Unlock()
    return fn
}

var (
    ErrMessageFormat   = fmt.Errorf("message format error")
    ErrHandlerNotFound = fmt.Errorf("message handler not found")
)

func NewRPCServer() *RPCServer {
    r := &RPCServer{}
    r.rpcMap.init()
    return r
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
    fn := s.rpcMap.Get(*msg.Name)
    if fn == nil {
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
    if s.onEvent != nil {
        s.onEvent(EventServe, addr)
    }
}
func (s *RPCServer) OnNew(client *Client) {
    if s.onEvent != nil {
        s.onEvent(EventAccept, client)
    }
}
func (s *RPCServer) OnClose(client *Client) {
    if s.onEvent != nil {
        s.onEvent(EventClose, client)
    }
}

func (s *RPCServer) HandlePacket(client *Client, packet *Packet) {
    Go(func() {
        if s.rateLimit != nil {
            s.rateLimit.Take()
        }
        payload := packet.Payload()
        msg := allocMessage()
        defer msg.Release()

        err := MSGUnpack(payload, msg)
        if err != nil {
            log.Println(err)
            return
        }
        s.onMessage(client, msg)
    })
}
func (s *RPCServer) AddPlugin(p interface{}) {
    if v, ok := p.(RateLimit); ok {
        s.rateLimit = v
    }
}
func (s *RPCServer) OnEvent(fn func(event RPCEvent, args ...interface{})) {
    s.onEvent = fn
}
func (s *RPCServer) RegisterFunc(name string, fn interface{}) {
    s.rpcMap.RegisterFunc(name, fn)
}
