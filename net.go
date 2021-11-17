package gof

import (
    "net"
    "time"
)

type (
    Delegate interface {
        OnStartServe(addr net.Addr)
        OnNew(conn net.Conn)
        OnClose(conn net.Conn)
        HandlePacket(packet *Packet)
    }
)
const (
    AutoRestartServerDuration = time.Second*3
)

