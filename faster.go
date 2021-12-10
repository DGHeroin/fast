package fast

import (
    "net"
)

type (
    ServerDelegate interface {
        OnStartServe(addr net.Addr)
        OnNew(client *Client)
        OnClose(client *Client)
        HandlePacket(client *Client, packet *Packet)
        AddPlugin(interface{})
    }
    ClientDelegate interface {
        OnOpen(client *Client)
        OnClose(client *Client)
        OnError(client *Client, err error)
        HandlePacket(client *Client, packet *Packet)
        AddPlugin(interface{})
    }
)

const (
//AutoRestartServerDuration = time.Second * 3
)

func Serve(network string, delegate ServerDelegate, opt Option) {
    RunForeverUntilPanic(opt.RetryDuration, func() {
        switch network {
        case "tcp":
            serveTCP(delegate, opt)
        case "kcp":
            serveKCP(delegate, opt)
        default:
            panic("unsupported protocol")
        }
    })
}
func NewClient(network string, delegate ClientDelegate, option Option) *Client {
    switch network {
    case "tcp":
        return newTCPClient(delegate, option)
    case "kcp":
        return newKCPClient(delegate, option)
    default:
        panic("unsupported protocol")
    }
}
