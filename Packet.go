package gf

import (
    "encoding/binary"
    "fmt"
    "sync"
    "unsafe"
)

const (
    MaxPayloadLength    = 64 * 1024 * 1024 // default 64Mb
    reserveHeaderSize   = 4
    preAllocPayloadSize = 128
)

var (
    ErrPayloadOverflow = fmt.Errorf("payload overflow")
    ErrPayloadTooLarge = fmt.Errorf("payload size too large")
)

var (
    packetPool = &sync.Pool{
        New: func() interface{} {
            p := &Packet{}
            return p
        },
    }
    payloadCapacities []uint32
    bufferPool        = map[uint32]*sync.Pool{}
)

func init() {
    // 预先分配好的各类byte buffer
    payloadCap := preAllocPayloadSize
    for payloadCap < MaxPayloadLength {
        payloadCapacities = append(payloadCapacities, uint32(payloadCap))
        payloadCap <<= 2
    }
    payloadCapacities = append(payloadCapacities, uint32(MaxPayloadLength))

    for _, capT := range payloadCapacities {
        bufferPool[capT] = &sync.Pool{
            New: func() interface{} {
                return make([]byte, reserveHeaderSize+capT)
            },
        }
    }
}

func allocPacket() *Packet {
    pkt := packetPool.Get().(*Packet)
    pkt.payload = make([]byte, reserveHeaderSize+preAllocPayloadSize)
    return pkt
}
func allocPayload(size uint32) ([]byte, *sync.Pool) {
    capT := uint32(0)
    for _, capacity := range payloadCapacities {
        if capacity > size {
            capT = capacity
        }
    }
    pool := bufferPool[capT]
    data := pool.Get().([]byte)
    return data, pool
}

type (
    Packet struct {
        payload    []byte
        readCursor int
    }
)

func NewPacket() *Packet {
    return allocPacket()
}

func (p *Packet) Release() {
    p.setLen(0)
    p.readCursor = 0
    packetPool.Put(p)
}

// read

func (p *Packet) ReadUint16() uint16 {
    return binary.LittleEndian.Uint16(p.ReadBytes(2))
}

// write
func (p *Packet) WriteUint16(v uint16) {
    pl := p.extendPayload(2)
    binary.LittleEndian.PutUint16(pl, v)
}
func (p *Packet) ReadBytes(size int) []byte {
    pos := p.readCursor
    end := pos + size
    if size > MaxPayloadLength || pos > int(p.Len()) {
        panic(ErrPayloadOverflow)
    }
    p.readCursor = end
    return p.payloadSlice(pos, end)
}
func (p *Packet) WriteBytes(b []byte) {
    payload := p.extendPayload(len(b))
    copy(payload, b)
}

func (p *Packet) Len() uint32 {
    binary.LittleEndian.Uint32(p.payload[0:4])
    return *(*uint32)(unsafe.Pointer(&p.payload[0]))
}
func (p *Packet) setLen(len uint32) {
    ptr := (*uint32)(unsafe.Pointer(&p.payload[0]))
    *ptr = len
}
func (p *Packet) Cap() uint32 {
    return uint32(len(p.payload) - reserveHeaderSize)
}

func (p *Packet) payloadSlice(start int, end int) []byte {
    return p.payload[reserveHeaderSize+start : reserveHeaderSize+end]
}
func (p *Packet) extendPayload(size int) []byte {
    if size > MaxPayloadLength || size < 0 {
        panic(ErrPayloadTooLarge)
    }

    oldCap := p.Cap()
    oldLen := p.Len()
    wantLen := oldLen + uint32(size)
    if oldCap >= wantLen { // 足够用
        p.setLen(wantLen)
        return p.payloadSlice(int(oldLen), int(wantLen))
    }

    oldPayload := p.payload
    newPayload, capPool := allocPayload(wantLen) // 重新申请buffer

    copy(newPayload, oldPayload)

    if oldCap > preAllocPayloadSize { // 归还buffer
        capPool.Put(oldPayload)
    }
    p.setLen(wantLen)
    return p.payloadSlice(int(oldLen), int(wantLen))
}
func (p *Packet) data() []byte {
    return p.payload[0 : reserveHeaderSize+p.Len()]
}
func (p *Packet) PayloadAsString() string {
    return string(p.Payload())
}

func (p *Packet) Payload() []byte {
    return p.payload[reserveHeaderSize : reserveHeaderSize+p.Len()]
}
