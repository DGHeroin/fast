package gf

import (
    "github.com/xtaci/kcp-go"
    "net"
)

func newKCPClient(delegate ClientDelegate, opt Option) *Client {
    var (
        conn net.Conn
        err  error
    )
    c := &Client{
        address: opt.Address,
    }
    openAndRead := func() {
        conn, err = kcp.Dial(opt.Address)
        if err != nil {
            return
        }
        c.Conn = conn
        pc := c.PacketConnection()

        defer func() {
            delegate.OnClose(c)
            pc.Close()
            c.packetConn = nil
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
