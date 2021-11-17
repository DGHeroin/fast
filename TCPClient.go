package gof

import (
    "crypto/tls"
    "net"
)

func newTCPClient(delegate ClientDelegate, opt Option) *Client {
    var (
        conn net.Conn
        err  error
    )
    c := &Client{
        address: opt.Address,
    }
    openAndRead := func() {
        if opt.TLS == nil {
            addr, err := net.ResolveTCPAddr("tcp", opt.Address)
            if err != nil {
                delegate.OnError(c, err)
                return
            }
            laddr, err := net.ResolveTCPAddr("tcp", ":0")
            if err != nil {
                delegate.OnError(c, err)
                return
            }
            conn, err = net.DialTCP("tcp", laddr, addr)
            if err != nil {
                delegate.OnError(c, err)
                return
            }
        } else {
            conn, err = tls.Dial("tcp", opt.Address, opt.TLS)
            delegate.OnError(c, err)
            return
        }
        c.Conn = conn
        pc := c.PacketConnection()

        defer func() {
            delegate.OnClose(c)
            pc.Close()
        }()

        go delegate.OnOpen(c)
        pc.LoopReadPack(func(packet *Packet, err error) {
            if err != nil {
                delegate.OnError(c, err)
                return
            }
            delegate.HandlePacket(c, packet)
        })
    }
    go RunForeverUntilPanic(opt.RetryDuration, openAndRead)

    return c
}
