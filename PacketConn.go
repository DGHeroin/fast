package gf

import (
    "encoding/binary"
    "net"
    "sync"
    "time"
)

type (
    PacketConn struct {
        conn       net.Conn
        sendChan   chan *Packet
        sendBuffer []*Packet
        closeCh    chan struct{}
        closeOnce sync.Once
    }
)

func (pc*PacketConn) LoopReadPack(cb func(packet *Packet, err error)) {
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
    pc.sendChan <- packet
}

func (pc *PacketConn) loop() {
    defer func() {
        recover()
    }()
    ticker := time.NewTicker(time.Millisecond * 5)
    defer ticker.Stop()
    for {
        select {
        case packet := <-pc.sendChan:
            pc.sendBuffer = append(pc.sendBuffer, packet)
        case <-ticker.C:
            copyBuffer := pc.sendBuffer
            pc.sendBuffer = []*Packet{}
            pc.flushSend(copyBuffer...)
        case <-pc.closeCh:
            return
        }
    }
}

func (pc *PacketConn) flushSend(packets ...*Packet) error {
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
    pc.closeOnce.Do(func() {
        close(pc.closeCh)
        pc.conn.Close()
    })
}
func NewPacketConnection(conn net.Conn) *PacketConn {
    pc := &PacketConn{
        conn:     NewConnection(conn, 1024*1024, 1024*1024),
        sendChan: make(chan *Packet),
        closeCh:  make(chan struct{}),
    }
    Go(pc.loop)
    return pc
}
