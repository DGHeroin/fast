package gf

import (
    "crypto/tls"
    "net"
    "time"
)

func serveTCP(delegate ServerDelegate, opt Option) {
    var (
        ln  net.Listener
        err error
    )
    if opt.TLS != nil {
        ln, err = tls.Listen("tcp", opt.Address, opt.TLS)
    } else {
        ln, err = net.Listen("tcp", opt.Address)
    }
    if err != nil {
        panic(err)
    }

    delegate.OnStartServe(ln.Addr())
    defer ln.Close()
    handleConn := func(conn net.Conn) {
        c := &Client{
            address: opt.Address,
        }
        c.Conn = conn
        pc := c.PacketConnection()
        defer func() {
            pc.Close()
            conn.Close()
            delegate.OnClose(c)
        }()

        delegate.OnNew(c)

        pc.LoopReadPack(func(packet *Packet, err error) {
            if err != nil {
                return
            }
            delegate.HandlePacket(c, packet)
        })
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
