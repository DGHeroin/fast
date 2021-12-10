package fast

import (
    "fmt"
    "net"
    "sync/atomic"
)

type (
    Client struct {
        id         int32
        address    string
        packetConn *PacketConn
        net.Conn
    }
)

var (
    rpcClientId = int32(0)
)

func newClient() *Client {
    cli := &Client{
        id: atomic.AddInt32(&rpcClientId, 1),
    }
    return cli
}

func (c *Client) SendPacket(data []byte) {
    if c.packetConn == nil {
        return
    }
    pc := c.packetConn
    packet := NewPacket()
    packet.WriteBytes(data)
    pc.Send(packet)
}

func (c *Client) PacketConnection() *PacketConn {
    if c.packetConn == nil {
        c.packetConn = NewPacketConnection(c.NetConn())
    }
    return c.packetConn
}

func (c *Client) NetConn() net.Conn {
    return c.Conn
}
func (c *Client) Close() {
    if c.Conn == nil {
        return
    }
    defer func() {
        recover()
    }()
    c.packetConn = nil
    c.Conn.Close()
    c.Conn = nil
}

func (c *Client) String() string {
    if c.Conn == nil {
        return fmt.Sprintf("Client<%s> closed", c.address)
    }
    return fmt.Sprintf("Client<%v> => <%v>", c.Conn.LocalAddr(), c.Conn.RemoteAddr())
}
