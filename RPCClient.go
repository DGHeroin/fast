package gf

import (
    "context"
    "fmt"
    "log"
    "sync"
    "sync/atomic"
    "time"
)

type (
    RPCClient struct {
        ClientHandler
        client   *Client
        mu       sync.Mutex
        handlers map[int32]*RPCRequestHandler
        state    int32
    }
    RPCRequestHandler struct {
        respPayload []byte
        OnDone      func()
        Err         []byte
    }
)

var (
    rpcId int32
)

func NewRPCClient() *RPCClient {
    c := &RPCClient{
        handlers: map[int32]*RPCRequestHandler{},
    }
    return c
}
func (c *RPCClient) Call(ctx context.Context, name string, r interface{}, w interface{}) (err error) {
    defer func() {
        if e := recover(); e != nil {
            err = fmt.Errorf("%v", e)
        }
    }()
    msg := newMessage(mTypeRequest)
    timeoutCh := time.After(time.Second * 120)

    defer msg.Release()
    msg.Id = atomic.AddInt32(&rpcId, 1)
    if r != nil {
        if payload, err := MSGPack(r); err != nil {
            return err
        } else {
            newName := name
            msg.Name = &newName
            msg.Payload = payload
        }
    }
    respChan := make(chan struct{})
    c.mu.Lock()

    var (
        handler  = &RPCRequestHandler{}
        doneOnce sync.Once
    )

    handler.OnDone = func() {
        doneOnce.Do(func() {
            close(respChan)
        })
    }
    defer handler.OnDone()
    c.handlers[msg.Id] = handler
    defer delete(c.handlers, msg.Id)
    c.mu.Unlock()

    sendData, _ := MSGPack(msg)

    c.client.SendPacket(sendData)

    select {
    case <-respChan: // 收到恢复
        if handler.Err != nil {
            return fmt.Errorf("%s", handler.Err)
        }
        return MSGUnpack(handler.respPayload, w)
    case <-ctx.Done(): // 取消?
        return fmt.Errorf("cancel")
    case <-timeoutCh:
        return fmt.Errorf("request timeout")
    }
}

func (c *RPCClient) onMessage(client *Client, msg *message) {
    switch msg.Type {
    case mTypeResponse:
        c.handleResponse(client, msg)
    }
}

func (c *RPCClient) handleResponse(client *Client, msg *message) {
    c.mu.Lock()
    h, ok := c.handlers[msg.Id]
    c.mu.Unlock()
    if !ok {
        return
    }
    h.respPayload = msg.Payload
    h.Err = msg.Err
    h.OnDone()

}

func (c *RPCClient) OnOpen(client *Client) {
    atomic.CompareAndSwapInt32(&c.state, 0, 1)
    log.Println("client open")
}
func (c *RPCClient) OnClose(client *Client) {
    atomic.CompareAndSwapInt32(&c.state, 1, 0)
    c.mu.Lock()
    for _, v := range c.handlers {
        v.Err = []byte("connection closed")
        v.OnDone()
    }
    c.handlers = map[int32]*RPCRequestHandler{}
    c.mu.Unlock()
}
func (c *RPCClient) OnError(client *Client, err error) {
    log.Printf("%v error:%v", client, err)
}
func (c *RPCClient) HandlePacket(client *Client, packet *Packet) {
    payload := packet.Payload()
    msg := allocMessage()
    err := MSGUnpack(payload, msg)
    if err != nil {
        log.Println(err)
        return
    }
    c.onMessage(client, msg)
}
func (c *RPCClient) Connect(network string, address string) {
    c.client = NewClient(network, c, Option{
        Address:       address,
        RetryDuration: time.Second,
    })
}
func (c *RPCClient) IsConnected() bool {
    return atomic.LoadInt32(&c.state) == 1
}
