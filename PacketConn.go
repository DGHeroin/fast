package fast

import (
    "encoding/binary"
    "net"
    "sync"
)

type (
    PacketConn struct {
        id         int32
        conn       net.Conn
        sendM      sync.Mutex
        sendBuffer []*Packet
        closeOnce  sync.Once
        rateLimit  RateLimit
    }
)

func (pc *PacketConn) LoopReadPack(cb func(packet *Packet, err error)) {
    for {
        pkt, err := pc.Recv()
        cb(pkt, err)
        if err != nil {
            break
        }
    }
}
func (pc *PacketConn) Recv() (*Packet, error) {
    var (
        header [4]byte
    )
    err := readFull(pc.conn, header[:])
    if err != nil {
        return nil, err
    }
    payloadSize := binary.LittleEndian.Uint32(header[:])
    if payloadSize > MaxPayloadLength {
        return nil, ErrPayloadTooLarge
    }
    packet := NewPacket()
    payload := packet.extendPayload(int(payloadSize))

    err = readFull(pc.conn, payload)
    if err != nil {
        packet.Release()
        return nil, err
    }
    packet.setLen(payloadSize)

    return packet, nil
}

func (pc *PacketConn) Send(packet *Packet) {
    pc.sendM.Lock()
    pc.sendBuffer = append(pc.sendBuffer, packet)
    pc.sendM.Unlock()

    if pc.rateLimit != nil {
        pc.rateLimit.Take()
    }
    Go(func() {
        pc.sendM.Lock()
        defer pc.sendM.Unlock()

        _ = pc.flushSend(pc.sendBuffer...)
        pc.sendBuffer = []*Packet{}
    })
}

func (pc *PacketConn) flushSend(packets ...*Packet) error {
    defer func() {
        recover()
    }()
    for _, packet := range packets {
        data := packet.data()
        _, err := pc.conn.Write(data)

        packet.Release()
        if err != nil {
            return err
        }
    }
    return tryFlush(pc.conn)
}

func (pc *PacketConn) Close() {
    _ = pc.conn.Close()
}
func (pc *PacketConn) AddPlugin(p interface{}) {
    if v, ok := p.(RateLimit); ok {
        pc.rateLimit = v
    }
}

func NewPacketConnection(conn net.Conn) *PacketConn {
    pc := &PacketConn{
        conn: NewConnection(conn, 1024*1024, 1024*1024),
    }

    return pc
}
