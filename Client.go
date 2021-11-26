package gf

import (
    "fmt"
    "net"
)

type Client struct {
    address    string
    packetConn *PacketConn
    net.Conn
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