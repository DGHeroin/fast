package gof

import "net"

type Client struct {
    address    string
    packetConn *PacketConn
    net.Conn
}

func (c *Client) SendPacket(data []byte) {
    pc := c.PacketConnection()
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

