package fast

import "sync"

type (
    message struct {
        Id      int32   `msgpack:"id"`
        Type    uint8   `msgpack:"type"`
        Name    *string `msgpack:"name,omitempty"`
        Err     []byte  `msgpack:"err,omitempty"`
        Payload []byte  `msgpack:"p,omitempty"`
    }
)

const (
    mTypeHeartBeat = uint8(0)
    mTypeRequest   = uint8(1)
    mTypeResponse  = uint8(2)
    mTypeSingleWay = uint8(3)
)

var (
    msgPool = &sync.Pool{
        New: func() interface{} {
            p := &message{}
            return p
        },
    }
)

func allocMessage() *message {
    m := msgPool.Get().(*message)

    return m
}
func newMessage(t uint8) *message {
    msg := allocMessage()
    msg.Type = t
    return msg
}
func (p *message) Release() {
    msgPool.Put(p)
}
