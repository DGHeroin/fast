package gof

import (
    "github.com/xtaci/kcp-go"
    "net"
    "time"
)

func serveKCP(delegate Delegate, opt Option) {
    var (
        ln  net.Listener
        err error
    )
    ln, err = kcp.Listen(opt.Address)
    if err != nil {
        panic(err)
    }

    delegate.OnStartServe(ln.Addr())
    defer ln.Close()
    handleConn := func(conn net.Conn) {
        delegate.OnNew(conn)
        defer delegate.OnClose(conn)
        pc := NewPacketConnection(conn)
        for {
            pkt, err := pc.Recv()
            if err != nil {
                break
            }
            if pkt != nil {
                delegate.HandlePacket(pkt)
            }
        }
    }

    for {
        conn, err := ln.Accept()
        if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
            time.Sleep(time.Second)
            continue
        }
        Go(func() {
            handleConn(conn)
        })
    }
}